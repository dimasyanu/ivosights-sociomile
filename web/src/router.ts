import {
  createMemoryHistory,
  createRouter,
  type RouteRecordRaw,
} from "vue-router"
import { type Menu } from "./data/admin-menus"

const mapMenu = (items: Menu[], basePath: string[]): RouteRecordRaw[] =>
  items.map((x) => {
    const paths = [...basePath, x.name]
    return {
      path: x.name,
      name: paths.join("."),
      component: x.component,
      children: x.children ? mapMenu(x.children, paths) : [],
      redirect: x.redirect ? { name: x.redirect } : undefined,
    }
  })

const routes = [
  { path: "/", redirect: { name: "login" } },
  {
    path: "/login",
    name: "login",
    component: () => import("@/pages/LoginPage.vue"),
  },
  {
    path: "/conversations",
    name: "conversations",
    component: () => import("@/pages/ConversationsPage.vue"),
  },
  {
    path: "/tickets",
    name: "tickets",
    component: () => import("@/pages/TicketsPage.vue"),
  },
  {
    path: "/users",
    name: "users",
    component: () => import("@/pages/UsersPage.vue"),
  },
  {
    path: "/:pathMatch(.*)*", // Matches all unmatched paths
    name: "notfound",
    component: () => import("@/pages/NotFoundPage.vue"),
  },
]

const router = createRouter({
  history: createMemoryHistory(),
  routes: routes,
})

export default router
