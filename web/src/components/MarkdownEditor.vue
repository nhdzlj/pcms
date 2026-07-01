<template>
  <div class="md-editor">
    <el-tabs v-model="activeTab" class="md-tabs">
      <el-tab-pane label="编辑" name="edit">
        <el-input
          v-model="localContent"
          type="textarea"
          :rows="20"
          :placeholder="placeholder"
          class="md-textarea"
          @input="handleInput"
        />
      </el-tab-pane>
      <el-tab-pane label="预览" name="preview">
        <div
          class="md-preview markdown-body"
          v-html="renderedContent"
        />
      </el-tab-pane>
      <el-tab-pane v-if="showSplit" label="分屏" name="split">
        <div class="md-split">
          <el-input
            v-model="localContent"
            type="textarea"
            :rows="20"
            :placeholder="placeholder"
            class="md-textarea"
            @input="handleInput"
          />
          <div class="md-split-divider" />
          <div
            class="md-preview markdown-body"
            v-html="renderedContent"
          />
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 工具栏 -->
    <div class="md-toolbar">
      <el-button
        v-for="item in toolbar"
        :key="item.label"
        size="small"
        @click="insertMarkdown(item.before, item.after)"
      >
        {{ item.label }}
      </el-button>
      <el-upload
        :action="uploadUrl"
        :headers="uploadHeaders"
        :show-file-list="false"
        :before-upload="beforeUpload"
        @success="handleUploadSuccess"
        style="display: inline-block; margin-left: 8px"
      >
        <el-button size="small">
          <el-icon><Picture /></el-icon> 图片
        </el-button>
      </el-upload>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { ElMessage } from "element-plus";
import { Picture } from "@element-plus/icons-vue";

const props = defineProps<{
  modelValue: string;
  placeholder?: string;
  showSplit?: boolean;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const activeTab = ref("edit");
const localContent = ref(props.modelValue || "");

const uploadUrl = "/api/v1/files/upload";
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem("token") || ""}`,
}));

watch(
  () => props.modelValue,
  (val) => {
    if (val !== localContent.value) {
      localContent.value = val || "";
    }
  }
);

function handleInput() {
  emit("update:modelValue", localContent.value);
}

interface ToolbarItem {
  label: string;
  before: string;
  after: string;
}

const toolbar: ToolbarItem[] = [
  { label: "H1", before: "# ", after: "" },
  { label: "H2", before: "## ", after: "" },
  { label: "H3", before: "### ", after: "" },
  { label: "B", before: "**", after: "**" },
  { label: "I", before: "*", after: "*" },
  { label: "~~", before: "~~", after: "~~" },
  { label: "`", before: "`", after: "`" },
  { label: "链接", before: "[", after: "](url)" },
  { label: "列表", before: "- ", after: "" },
  { label: "引用", before: "> ", after: "" },
  { label: "代码块", before: "```\n", after: "\n```" },
  { label: "表格", before: "", after: "\n| col1 | col2 | col3 |\n| --- | --- | --- |\n| | | |\n" },
];

function insertMarkdown(before: string, after: string) {
  const textarea = document.querySelector(".md-textarea textarea") as HTMLTextAreaElement;
  if (!textarea) return;

  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const selectedText = localContent.value.slice(start, end);

  const newText =
    localContent.value.slice(0, start) +
    before +
    selectedText +
    after +
    localContent.value.slice(end);

  localContent.value = newText;
  emit("update:modelValue", newText);

  // Restore cursor position
  nextTick(() => {
    textarea.focus();
    textarea.setSelectionRange(start + before.length, start + before.length + selectedText.length);
  });
}

import { nextTick } from "vue";

function simpleMarkdown(text: string): string {
  if (!text) return "";
  let html = text;
  // Escape HTML
  html = html.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  // Code blocks
  html = html.replace(/```(\w*)\n([\s\S]*?)```/g, '<pre><code>$2</code></pre>');
  // Inline code
  html = html.replace(/`([^`]+)`/g, "<code>$1</code>");
  // Headers
  html = html.replace(/^#### (.+)$/gm, "<h4>$1</h4>");
  html = html.replace(/^### (.+)$/gm, "<h3>$1</h3>");
  html = html.replace(/^## (.+)$/gm, "<h2>$1</h2>");
  html = html.replace(/^# (.+)$/gm, "<h1>$1</h1>");
  // Bold & italic
  html = html.replace(/\*\*\*(.+?)\*\*\*/g, "<strong><em>$1</em></strong>");
  html = html.replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>");
  html = html.replace(/\*(.+?)\*/g, "<em>$1</em>");
  // Strikethrough
  html = html.replace(/~~(.+?)~~/g, "<del>$1</del>");
  // Images
  html = html.replace(/!\[([^\]]*)\]\(([^)]+)\)/g, '<img src="$2" alt="$1" />');
  // Links
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank">$1</a>');
  // Horizontal rule
  html = html.replace(/^---$/gm, "<hr />");
  // Blockquotes
  html = html.replace(/^> (.+)$/gm, "<blockquote>$1</blockquote>");
  // Unordered lists
  html = html.replace(/^[\-\*] (.+)$/gm, "<li>$1</li>");
  html = html.replace(/(<li>.*<\/li>)/gs, "<ul>$1</ul>");
  // Ordered lists
  html = html.replace(/^\d+\. (.+)$/gm, "<li>$1</li>");
  // Tables (simple)
  html = html.replace(/^\|(.+)\|$/gm, (match) => {
    const cells = match.split("|").filter((c) => c.trim());
    const isHeader = cells.some((c) => /^[\s\-:]+$/.test(c.trim()));
    if (isHeader) return "";
    const tag = match.includes("---") ? "" : "td";
    if (!tag) return "";
    return "<tr>" + cells.map((c) => `<${tag}>${c.trim()}</${tag}>`).join("") + "</tr>";
  });
  // Paragraphs
  html = html.replace(/\n\n+/g, "</p><p>");
  html = "<p>" + html + "</p>";
  // Clean empty
  html = html.replace(/<p><\/p>/g, "");
  // Fix consecutive blockquotes
  html = html.replace(/<\/blockquote>\s*<blockquote>/g, "<br />");

  return html;
}

const renderedContent = computed(() => simpleMarkdown(localContent.value));

function beforeUpload(file: File) {
  const isImage = file.type.startsWith("image/");
  if (!isImage) {
    ElMessage.warning("只能上传图片文件");
    return false;
  }
  const isLt10M = file.size / 1024 / 1024 < 10;
  if (!isLt10M) {
    ElMessage.warning("图片大小不能超过 10MB");
    return false;
  }
  return true;
}

function handleUploadSuccess(response: any) {
  if (response.data?.url) {
    insertMarkdown(`![image](${response.data.url})`, "");
    ElMessage.success("上传成功");
  }
}
</script>

<style scoped>
.md-editor {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.md-tabs {
  flex: 1;
  overflow: hidden;
}

.md-tabs :deep(.el-tabs__content) {
  height: calc(100% - 40px);
  overflow-y: auto;
}

.md-tabs :deep(.el-tab-pane) {
  height: 100%;
}

.md-textarea :deep(textarea) {
  font-family: "JetBrains Mono", "Fira Code", monospace;
  font-size: 14px;
  line-height: 1.8;
  height: 100% !important;
  border: none;
  resize: none;
  padding: 16px;
}

.md-preview {
  padding: 16px 24px;
  line-height: 1.8;
  font-size: 15px;
}

.md-preview :deep(h1) { font-size: 28px; margin: 8px 0 16px; border-bottom: 2px solid #eee; padding-bottom: 8px; }
.md-preview :deep(h2) { font-size: 24px; margin: 8px 0 12px; border-bottom: 1px solid #eee; padding-bottom: 6px; }
.md-preview :deep(h3) { font-size: 20px; margin: 8px 0 10px; }
.md-preview :deep(h4) { font-size: 16px; margin: 8px 0 8px; }
.md-preview :deep(p) { margin: 8px 0; }
.md-preview :deep(code) { background: #f4f4f5; padding: 2px 6px; border-radius: 4px; font-size: 13px; }
.md-preview :deep(pre) { background: #282c34; color: #abb2bf; padding: 16px; border-radius: 8px; overflow-x: auto; margin: 12px 0; }
.md-preview :deep(pre code) { background: transparent; padding: 0; color: inherit; }
.md-preview :deep(blockquote) { border-left: 4px solid var(--app-primary); padding: 4px 16px; margin: 12px 0; background: var(--app-primary-light); border-radius: 0 8px 8px 0; }
.md-preview :deep(img) { max-width: 100%; border-radius: 8px; }
.md-preview :deep(a) { color: var(--app-primary); }
.md-preview :deep(ul), .md-preview :deep(ol) { padding-left: 24px; margin: 8px 0; }
.md-preview :deep(li) { margin: 4px 0; }
.md-preview :deep(hr) { border: none; border-top: 1px solid #eee; margin: 20px 0; }
.md-preview :deep(table) { border-collapse: collapse; width: 100%; margin: 12px 0; }
.md-preview :deep(th), .md-preview :deep(td) { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }
.md-preview :deep(th) { background: #f5f7fa; font-weight: 600; }

.md-split {
  display: flex;
  height: 100%;
}

.md-split .md-textarea {
  flex: 1;
}

.md-split-divider {
  width: 4px;
  background: var(--app-border);
  cursor: col-resize;
}

.md-split .md-preview {
  flex: 1;
  overflow-y: auto;
}

.md-toolbar {
  padding: 8px 16px;
  border-top: 1px solid var(--app-border);
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
  background: var(--app-header-bg);
  flex-shrink: 0;
}
</style>
