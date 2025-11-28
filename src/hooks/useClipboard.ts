import { readFile, stat } from "@tauri-apps/plugin-fs";
import { useMount } from "ahooks";
import { cloneDeep } from "es-toolkit";
import { isEmpty, remove } from "es-toolkit/compat";
import { nanoid } from "nanoid";
import {
  type ClipboardChangeOptions,
  onClipboardChange,
  startListening,
} from "tauri-plugin-clipboard-x-api";
import { fullName } from "tauri-plugin-fs-pro-api";
import {
  insertHistory,
  selectHistory,
  updateHistory,
} from "@/database/history";
import type { State } from "@/pages/Main";
import { getClipboardTextSubtype } from "@/plugins/clipboard";
import { clipboardStore } from "@/stores/clipboard";
import { globalStore } from "@/stores/global";
import type { DatabaseSchemaHistory } from "@/types/database";
import { formatDate } from "@/utils/dayjs";
import { sendWebhook } from "@/utils/webhook";

export const useClipboard = (
  state: State,
  options?: ClipboardChangeOptions,
) => {
  useMount(async () => {
    await startListening();

    onClipboardChange(async (result) => {
      const { files, image, html, rtf, text } = result;

      if (isEmpty(result) || Object.values(result).every(isEmpty)) return;

      const { copyPlain } = clipboardStore.content;

      const data = {
        createTime: formatDate(),
        favorite: false,
        group: "text",
        id: nanoid(),
        search: text?.value,
      } as DatabaseSchemaHistory;

      if (files) {
        Object.assign(data, files, {
          group: "files",
          search: files.value.join(" "),
        });
      } else if (html && !copyPlain) {
        Object.assign(data, html);
      } else if (rtf && !copyPlain) {
        Object.assign(data, rtf);
      } else if (text) {
        const subtype = await getClipboardTextSubtype(text.value);

        Object.assign(data, text, {
          subtype,
        });
      } else if (image) {
        Object.assign(data, image, {
          group: "image",
        });
      }

      const sqlData = cloneDeep(data);

      const { type, value, group, createTime } = data;

      if (type === "image") {
        sqlData.value = await fullName(value);
      }

      if (type === "files") {
        sqlData.value = JSON.stringify(value);
      }

      const [matched] = await selectHistory((qb) => {
        const { type, value } = sqlData;

        return qb.where("type", "=", type).where("value", "=", value);
      });

      const visible = state.group === "all" || state.group === group;

      if (matched) {
        if (!clipboardStore.content.autoSort) return;

        const { id } = matched;

        if (visible) {
          remove(state.list, { id });

          state.list.unshift({ ...data, id });
        }

        return updateHistory(id, { createTime });
      }

      if (visible) {
        state.list.unshift(data);
      }

      insertHistory(sqlData);

      // 发送 Webhook
      if (globalStore.webhook.enabled && globalStore.webhook.url) {
        try {
          // 检查是否为图片文件且小于 2MB
          if (type === "files" && Array.isArray(value)) {
            for (const path of value) {
              if (/\.(jpg|jpeg|png|gif|webp|bmp|svg)$/i.test(path)) {
                try {
                  const fileStat = await stat(path);
                  if (fileStat.size < 2 * 1024 * 1024) {
                    // 2MB
                    const bytes = await readFile(path);
                    let binary = "";
                    const len = bytes.byteLength;
                    for (let i = 0; i < len; i++) {
                      binary += String.fromCharCode(bytes[i]);
                    }
                    const base64 = window.btoa(binary);
                    const ext = path.split(".").pop()?.toLowerCase();
                    const mime =
                      ext === "jpg" ? "jpeg" : ext === "svg" ? "svg+xml" : ext;
                    const dataUrl = `data:image/${mime};base64,${base64}`;

                    await sendWebhook(
                      "image",
                      dataUrl,
                      globalStore.webhook.url,
                      "base64",
                    );
                    break;
                  }
                } catch (err) {
                  // biome-ignore lint/suspicious/noConsole: Log error for debugging
                  console.error(
                    "Failed to process image file for webhook:",
                    err,
                  );
                }
              }
            }
          } else {
            try {
              const webhookValue =
                type === "files"
                  ? value
                  : type === "image"
                    ? sqlData.value
                    : value;

              await sendWebhook(
                group as "text" | "image" | "files" | "html" | "rtf",
                webhookValue,
                globalStore.webhook.url,
                data.subtype,
              );
            } catch (error) {
              // biome-ignore lint/suspicious/noConsole: Log error for debugging
              console.error("Failed to send webhook:", error);
            }
          }
        } catch (error) {
          // biome-ignore lint/suspicious/noConsole: Log error for debugging
          console.error("Failed to send webhook:", error);
        }
      }
    }, options);
  });
};
