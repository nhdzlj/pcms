import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/LoginView.vue"),
    meta: { title: "登录", noAuth: true },
  },
  {
    path: "/",
    name: "Layout",
    component: () => import("@/views/LayoutView.vue"),
    redirect: "/home",
    children: [
      {
        path: "home",
        name: "Home",
        component: () => import("@/views/HomeView.vue"),
        meta: { title: "首页" },
      },
      {
        path: "documents",
        name: "DocumentList",
        component: () => import("@/views/DocumentList.vue"),
        meta: { title: "文档列表" },
      },
      {
        path: "documents/new",
        name: "DocumentCreate",
        component: () => import("@/views/DocumentEditor.vue"),
        meta: { title: "新建文档" },
      },
      {
        path: "documents/:id",
        name: "DocumentEdit",
        component: () => import("@/views/DocumentEditor.vue"),
        meta: { title: "编辑文档" },
      },
      {
        path: "search",
        name: "Search",
        component: () => import("@/views/SearchView.vue"),
        meta: { title: "搜索" },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem("token");
  const requireAuth = !to.meta.noAuth;

  if (requireAuth && !token) {
    next("/login");
  } else if (to.path === "/login" && token) {
    next("/");
  } else {
    next();
  }
});

export default router;
