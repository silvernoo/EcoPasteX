<template>
  <div id="app" class="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900">
    <div class="container mx-auto px-4 py-8">
      <!-- Header -->
      <header class="mb-8">
        <h1 class="text-4xl font-bold text-white mb-2 flex items-center gap-3">
          <svg class="w-10 h-10 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          EcoPaste 剪贴板历史
        </h1>
        <p class="text-gray-400">管理和查看您的剪贴板历史记录</p>
      </header>

      <!-- Search and Filter Bar -->
      <div class="bg-gray-800 rounded-xl shadow-2xl p-6 mb-6 border border-gray-700">
        <div class="flex flex-col md:flex-row gap-4">
          <!-- Search Input -->
          <div class="flex-1">
            <div class="relative">
              <input
                v-model="searchQuery"
                @input="debouncedSearch"
                type="text"
                placeholder="搜索剪贴板内容..."
                class="w-full px-4 py-3 pl-12 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
              />
              <svg class="absolute left-4 top-3.5 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
          </div>

          <!-- Filter Buttons -->
          <div class="flex gap-2">
            <button
              v-for="filter in filters"
              :key="filter.value"
              @click="selectedFilter = filter.value; fetchItems()"
              :class="[
                'px-6 py-3 rounded-lg font-medium transition-all duration-200',
                selectedFilter === filter.value
                  ? 'bg-blue-600 text-white shadow-lg shadow-blue-500/50'
                  : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
              ]"
            >
              {{ filter.label }}
            </button>
          </div>
        </div>

        <!-- Stats -->
        <div class="mt-4 flex items-center gap-6 text-sm text-gray-400">
          <span>总计: <span class="text-white font-semibold">{{ totalItems }}</span> 条</span>
          <span>当前页: <span class="text-white font-semibold">{{ currentPage }}/{{ totalPages }}</span></span>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center py-20">
        <div class="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-blue-500"></div>
      </div>

      <!-- Items Grid -->
      <div v-else-if="items.length > 0" class="grid grid-cols-1 gap-4 mb-6">
        <div
          v-for="item in items"
          :key="item.id"
          class="bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-700 hover:border-blue-500 transition-all duration-200 hover:shadow-2xl hover:shadow-blue-500/20"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1 min-w-0">
              <!-- Type Badge -->
              <div class="flex items-center gap-2 mb-3">
                <span
                  :class="[
                    'px-3 py-1 rounded-full text-xs font-semibold',
                    item.isImage
                      ? 'bg-purple-500/20 text-purple-300 border border-purple-500/30'
                      : 'bg-green-500/20 text-green-300 border border-green-500/30'
                  ]"
                >
                  {{ item.isImage ? '图片' : '文本' }}
                </span>
                <span class="text-xs text-gray-500">
                  {{ formatDate(item.timestamp) }}
                </span>
              </div>

              <!-- Content Preview -->
              <div v-if="item.isImage" class="mb-3">
                <img
                  v-if="getImageSrc(item.value)"
                  :src="getImageSrc(item.value)"
                  alt="Clipboard Image"
                  class="max-w-full h-auto max-h-96 rounded-lg border border-gray-600"
                  @error="handleImageError"
                />
                <div v-else class="bg-gray-700 rounded-lg p-8 text-center">
                  <svg class="w-16 h-16 mx-auto text-gray-500 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  <p class="text-gray-400">图片预览不可用</p>
                </div>
              </div>
              <div v-else class="text-gray-300 leading-relaxed">
                <div v-if="item.preview" class="whitespace-pre-wrap break-words font-mono text-sm bg-gray-900 p-4 rounded-lg border border-gray-700">
                  {{ item.preview }}
                </div>
                <div v-else class="text-gray-500 italic">无预览内容</div>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-2">
              <!-- Copy Button -->
              <button
                @click="copyToClipboard(item)"
                class="flex-shrink-0 p-2 text-blue-400 hover:text-blue-300 hover:bg-blue-500/10 rounded-lg transition"
                title="复制"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                </svg>
              </button>

              <!-- Delete Button -->
            <button
              @click="deleteItem(item.id)"
              class="flex-shrink-0 p-2 text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-lg transition"
              title="删除"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-20">
        <svg class="w-24 h-24 mx-auto text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="text-xl font-semibold text-gray-400 mb-2">暂无数据</h3>
        <p class="text-gray-500">开始使用 EcoPaste 复制内容，数据将自动同步到这里</p>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center items-center gap-2">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage === 1"
          class="px-4 py-2 bg-gray-700 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition"
        >
          上一页
        </button>
        
        <div class="flex gap-1">
          <button
            v-for="page in visiblePages"
            :key="page"
            @click="goToPage(page)"
            :class="[
              'px-4 py-2 rounded-lg transition',
              currentPage === page
                ? 'bg-blue-600 text-white'
                : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
            ]"
          >
            {{ page }}
          </button>
        </div>

        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="px-4 py-2 bg-gray-700 text-white rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

const API_BASE_URL = 'api'

const items = ref([])
const loading = ref(false)
const searchQuery = ref('')
const selectedFilter = ref('all')
const currentPage = ref(1)
const pageSize = ref(10)
const totalItems = ref(0)
const totalPages = ref(0)

const filters = [
  { label: '全部', value: 'all' },
  { label: '文本', value: 'text' },
  { label: '图片', value: 'image' }
]

let searchTimeout = null

const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    fetchItems()
  }, 500)
}

const fetchItems = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      type: selectedFilter.value,
      search: searchQuery.value
    }

    const response = await axios.get(`${API_BASE_URL}/clipboard`, { params })
    items.value = response.data.items || []
    totalItems.value = response.data.total
    totalPages.value = response.data.totalPages
  } catch (error) {
    console.error('Failed to fetch items:', error)
  } finally {
    loading.value = false
  }
}

const deleteItem = async (id) => {
  if (!confirm('确定要删除这条记录吗？')) return

  try {
    await axios.delete(`${API_BASE_URL}/clipboard/${id}`)
    fetchItems()
  } catch (error) {
    console.error('Failed to delete item:', error)
    alert('删除失败，请重试')
  }
}

const cleanContent = (value) => {
  if (typeof value !== 'string') return value
  
  // 如果包含图片标签，提取 src
  if (value.includes('<img')) {
    const match = value.match(/src="([^"]+)"/)
    if (match) return match[1]
  }

  // 如果包含 HTML 标签，提取纯文本
  if (value.includes('<') && value.includes('>')) {
    const div = document.createElement('div')
    div.innerHTML = value
    return div.textContent || div.innerText || value
  }

  return value
}

const copyToClipboard = async (item) => {
  try {
    let contentToCopy = item.value
    if (typeof item.value === 'string') {
      contentToCopy = cleanContent(item.value)
    }

    if (item.isImage && typeof contentToCopy === 'string' && contentToCopy.startsWith('data:image')) {
      // 复制图片
      const response = await fetch(contentToCopy)
      const blob = await response.blob()
      await navigator.clipboard.write([
        new ClipboardItem({ [blob.type]: blob })
      ])
    } else {
      // 复制文本
      const text = typeof contentToCopy === 'string' ? contentToCopy : JSON.stringify(contentToCopy)
      await navigator.clipboard.writeText(text)
    }
    // 可以添加一个临时的成功提示，这里简单用 alert 或者后续集成 toast
    // alert('已复制') 
  } catch (err) {
    console.error('Failed to copy:', err)
    alert('复制失败: ' + err.message)
  }
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchItems()
  }
}

const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)

  if (end - start < maxVisible - 1) {
    start = Math.max(1, end - maxVisible + 1)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const formatDate = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} 天前`

  return date.toLocaleString('zh-CN')
}

const getImageSrc = (value) => {
  if (!value) return null
  if (typeof value === 'string') {
    // 处理 base64 图片
    if (value.startsWith('data:image')) {
      return value
    }
    // 处理直接的图片 URL
    if (value.startsWith('http') && (value.match(/\.(jpeg|jpg|gif|png|webp|bmp|svg)/i) || value.includes('images'))) {
      return value
    }
    // 处理 HTML 中的图片
    if (value.includes('<img')) {
      const match = value.match(/src="([^"]+)"/)
      if (match) return match[1]
    }
  }
  return null
}

const handleImageError = (e) => {
  e.target.style.display = 'none'
}

onMounted(() => {
  fetchItems()
  // 每30秒自动刷新
  setInterval(fetchItems, 30000)
})
</script>

<style scoped>
/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #1f2937;
}

::-webkit-scrollbar-thumb {
  background: #4b5563;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #6b7280;
}
</style>
