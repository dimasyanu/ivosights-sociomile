import type { RouteComponent } from 'vue-router'

type Lazy<T> = () => Promise<T>

export type Menu = {
  isReadOnly?: boolean
  label: string
  name: string
  alias?: string
  icon?: string
  component?: Lazy<RouteComponent>
  children?: Menu[]
  redirect?: string
}
const menus: Menu[] = [
  { isReadOnly: true, label: 'Overview', name: 'overview' },
  {
    label: 'Conversations',
    name: 'conversations',
    icon: 'conversation',
    component: () => import('@/pages/ConversationsPage.vue'),
  },
  {
    label: 'Tickets',
    name: 'tickets',
    icon: 'tickets',
    component: () => import('@/pages/TicketsPage.vue'),
  },
  {
    label: 'Users',
    name: 'users',
    icon: 'users',
    component: () => import('@/pages/UsersPage.vue'),
  },
]

export const adminMenus = menus
