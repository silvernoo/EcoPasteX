interface WebhookPayload {
  type: "text" | "image" | "files" | "html" | "rtf";
  value: string | string[];
  timestamp: string;
  subtype?: string;
}

interface QueueItem {
  payload: WebhookPayload;
  retries: number;
}

class WebhookQueue {
  private queue: QueueItem[] = [];
  private isProcessing = false;
  private maxRetries = 3;
  private retryDelay = 1000; // 1 second

  async add(payload: WebhookPayload, url: string) {
    this.queue.push({ payload, retries: 0 });

    if (!this.isProcessing) {
      await this.process(url);
    }
  }

  private async process(url: string) {
    if (this.queue.length === 0) {
      this.isProcessing = false;
      return;
    }

    this.isProcessing = true;
    const item = this.queue[0];

    try {
      await this.sendRequest(url, item.payload);
      // 成功后移除队列项
      this.queue.shift();
    } catch (error) {
      // biome-ignore lint/suspicious/noConsole: Log error for debugging
      console.error("Webhook request failed:", error);

      item.retries++;

      if (item.retries >= this.maxRetries) {
        // biome-ignore lint/suspicious/noConsole: Log error for debugging
        console.error("Max retries reached, removing item from queue");
        this.queue.shift();
      } else {
        // 重试前等待
        await new Promise((resolve) => setTimeout(resolve, this.retryDelay));
      }
    }

    // 继续处理下一个
    await this.process(url);
  }

  private async sendRequest(url: string, payload: WebhookPayload) {
    if (!url) {
      throw new Error("Webhook URL is not configured");
    }

    const response = await fetch(url, {
      body: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
      method: "POST",
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response;
  }

  getQueueLength(): number {
    return this.queue.length;
  }

  clear() {
    this.queue = [];
    this.isProcessing = false;
  }
}

export const webhookQueue = new WebhookQueue();

export const sendWebhook = async (
  type: WebhookPayload["type"],
  value: string | string[],
  url: string,
  subtype?: string,
) => {
  const payload: WebhookPayload = {
    subtype,
    timestamp: new Date().toISOString(),
    type,
    value,
  };

  await webhookQueue.add(payload, url);
};
