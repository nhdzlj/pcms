<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-title">PCMS</div>
      <div class="login-subtitle">个人认知管理系统</div>

      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleSubmit">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            show-password
            :prefix-icon="Lock"
            @keyup.enter="handleSubmit"
          />
        </el-form-item>
        <el-form-item v-if="isRegister" prop="email">
          <el-input
            v-model="form.email"
            placeholder="邮箱（选填）"
            size="large"
            :prefix-icon="Message"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            style="width: 100%"
            @click="handleSubmit"
          >
            {{ isRegister ? "注册" : "登录" }}
          </el-button>
        </el-form-item>
      </el-form>

      <div style="text-align: center">
        <el-button link type="primary" @click="isRegister = !isRegister">
          {{ isRegister ? "已有账号？去登录" : "没有账号？去注册" }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { User, Lock, Message } from "@element-plus/icons-vue";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const authStore = useAuthStore();

const isRegister = ref(false);
const loading = ref(false);
const formRef = ref();

const form = reactive({
  username: "",
  password: "",
  email: "",
});

const rules = {
  username: [
    { required: true, message: "请输入用户名", trigger: "blur" },
    { min: 2, max: 64, message: "用户名长度 2-64", trigger: "blur" },
  ],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, max: 64, message: "密码长度 6-64", trigger: "blur" },
  ],
};

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  loading.value = true;
  try {
    if (isRegister.value) {
      await authStore.register(form.username, form.password, form.email);
      ElMessage.success("注册成功");
    } else {
      await authStore.login(form.username, form.password);
      ElMessage.success("登录成功");
    }
    router.push("/");
  } catch (err: any) {
    ElMessage.error(err.message || "操作失败");
  } finally {
    loading.value = false;
  }
}
</script>
