package main

import (
	"context"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClipboardItem struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type      string             `json:"type" bson:"type"`
	Value     interface{}        `json:"value" bson:"value"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Subtype   string             `json:"subtype,omitempty" bson:"subtype,omitempty"`
	IsImage   bool               `json:"isImage" bson:"isImage"`
	Preview   string             `json:"preview,omitempty" bson:"preview,omitempty"`
}

type WebhookPayload struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Timestamp string      `json:"timestamp"`
	Subtype   string      `json:"subtype,omitempty"`
}

type PaginationResponse struct {
	Items      []ClipboardItem `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
	TotalPages int             `json:"totalPages"`
}

var (
	mongoClient *mongo.Client
	collection  *mongo.Collection
)

func initMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient = client
	collection = client.Database("ecopaste").Collection("clipboard")

	// 创建索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: -1},
		},
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Failed to create index: %v", err)
	}

	log.Println("Connected to MongoDB!")
}

func handleWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timestamp, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}

	item := ClipboardItem{
		Type:      payload.Type,
		Value:     payload.Value,
		Timestamp: timestamp,
		Subtype:   payload.Subtype,
		IsImage:   payload.Type == "image",
	}

	// 处理图片类型，提取预览
	if payload.Type == "image" {
		if valueStr, ok := payload.Value.(string); ok {
			item.Preview = extractImagePreview(valueStr)
		}
	} else if payload.Type == "text" || payload.Type == "html" {
		if valueStr, ok := payload.Value.(string); ok {
			// 检查是否包含 HTML 图片标签
			if strings.Contains(valueStr, "<img") && strings.Contains(valueStr, "src=") {
				item.IsImage = true
				item.Preview = "Image"
			} else if isImageURL(valueStr) {
				item.IsImage = true
				item.Preview = "Image (URL)"
			} else {
				// 处理文本类型，提取预览
				item.Preview = extractTextPreview(valueStr)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data received successfully"})
}

func getClipboardItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	filterType := c.Query("type")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 100
	}

	skip := (page - 1) * pageSize

	// 构建过滤条件
	filter := bson.M{}
	if filterType != "" && filterType != "all" {
		if filterType == "image" {
			filter["isImage"] = true
		} else if filterType == "text" {
			filter["isImage"] = false
		}
	}

	if search != "" {
		filter["preview"] = bson.M{"$regex": search, "$options": "i"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count documents"})
		return
	}

	// 查询数据
	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer cursor.Close(ctx)

	var items []ClipboardItem
	if err = cursor.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode data"})
		return
	}

	if items == nil {
		items = []ClipboardItem{}
	}

	// 修正现有数据的 IsImage 字段
	for i := range items {
		if !items[i].IsImage {
			if str, ok := items[i].Value.(string); ok {
				if strings.Contains(str, "<img") && strings.Contains(str, "src=") {
					items[i].IsImage = true
				} else if isImageURL(str) {
					items[i].IsImage = true
				}
			}
		}
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := PaginationResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

func isImageURL(s string) bool {
	if !strings.HasPrefix(s, "http") {
		return false
	}
	lower := strings.ToLower(s)
	return strings.Contains(lower, ".jpg") ||
		strings.Contains(lower, ".jpeg") ||
		strings.Contains(lower, ".png") ||
		strings.Contains(lower, ".gif") ||
		strings.Contains(lower, ".webp") ||
		strings.Contains(lower, ".bmp") ||
		strings.Contains(lower, ".svg")
}

func extractImagePreview(value string) string {
	// 如果是 base64 图片，提取前100个字符作为预览标识
	if strings.Contains(value, "base64,") {
		return "Image (base64)"
	}
	return "Image"
}

func extractTextPreview(value string) string {
	// 移除 HTML 标签
	text := stripHTMLTags(value)

	// 限制长度
	maxLen := 200
	if len(text) > maxLen {
		return text[:maxLen] + "..."
	}
	return text
}

func stripHTMLTags(html string) string {
	// 简单的 HTML 标签移除
	text := html
	for {
		start := strings.Index(text, "<")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}
	return strings.TrimSpace(text)
}

func deleteItem(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

func main() {
	initMongoDB()
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Frontend static files serving
	frontendDistPath := path.Join("dist")
	r.Static("/assets", path.Join(frontendDistPath, "assets"))               // Serve assets from /assets
	r.StaticFile("/", path.Join(frontendDistPath, "index.html"))             // Serve index.html for root
	r.StaticFile("/favicon.ico", path.Join(frontendDistPath, "favicon.ico")) // Serve favicon

	r.NoRoute(func(c *gin.Context) {
		// If the path does not contain "/api", serve index.html
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.File(path.Join(frontendDistPath, "index.html"))
		}
	})

	// API 路由
	api := r.Group("/api")
	{
		api.POST("/webhook", handleWebhook)
		api.GET("/clipboard", getClipboardItems)
		api.DELETE("/clipboard/:id", deleteItem)
	}
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("Server starting on :3000")
	r.Run(":3000")
}
