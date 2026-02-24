<script setup lang="ts">
import SvgIcon from "@/components/SvgIcon.vue"
import ThemeToggler from "@/components/ThemeToggler.vue"
import { useGlobalStore } from "@/stores/global"
import { computed, ref } from "vue"
import { useRouter } from "vue-router"
import auth from "@/api/server/auth"

const email = ref("")
const password = ref("")

const router = useRouter()
const globalStore = useGlobalStore()
const isPasswordShown = ref(false)
const togglePasswordVisibility = (e: MouseEvent) => {
  e.preventDefault()
  isPasswordShown.value = !isPasswordShown.value
}
const fgColor = computed(() =>
  globalStore.theme === "light" ? "#333" : "#fff",
)
const handleSubmit = async () => {
  try {
    await auth.login(email.value, password.value)
    router.push({ name: "conversations" })
  } catch (error) {
    console.error("Login failed:", error)
  }
}
</script>

<template>
  <main>
    <div class="banner"></div>
    <div class="form-container">
      <div class="heading flex justify-between items-center">
        <a href="/">
          <img src="/app.svg" alt="App Logo" class="h-16 w-16" />
        </a>
        <ThemeToggler />
      </div>
      <div class="container flex h-max m-auto justify-center items-center">
        <form @submit.prevent="handleSubmit">
          <h3 class="text-xl font-semibold text-center mt-[1rem]">Login</h3>

          <!-- Email Address input -->
          <div class="fieldset w-full">
            <legend class="fieldset-legend">Email Address</legend>
            <label class="input" for="email">
              <SvgIcon name="mail" :stroke="fgColor" :width="32" />
              <input id="email" type="email" class="grow pl-1" name="email" placeholder="Email Address"
                autocomplete="email" v-model="email" />
            </label>
          </div>

          <!-- Password input -->
          <div class="fieldset w-full">
            <legend class="fieldset-legend">Password</legend>
            <label class="input" for="email">
              <SvgIcon name="key" :stroke="fgColor" :width="32" />
              <input id="password" name="password" :type="isPasswordShown ? 'text' : 'password'" class="grow pl-1"
                placeholder="Password" v-model="password" />
              <button class="cursor-pointer" @click="togglePasswordVisibility($event)">
                <SvgIcon name="eye" v-if="isPasswordShown" :stroke="fgColor" />
                <SvgIcon name="eye-slash" v-else :stroke="fgColor" />
              </button>
            </label>
          </div>

          <!-- Login button -->
          <div class="fieldset mt-4">
            <button class="btn btn-primary">
              <SvgIcon name="arrow-down-tray" class="-rotate-90" stroke="#ffffff" />
              <span class="ml-2">Login</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </main>
</template>

<style scoped lang="scss">
@use "../sass/size.scss" as size;

main {
  display: flex;
  width: 100%;
  height: 100vh;
  flex-direction: row;

  &>.banner {
    flex: 1 1 auto;
    background-color: var(--color-base-200);

    display: none;

    @media (min-width: size.$smMin) {
      display: flex;
    }
  }

  &>.form-container {
    flex: 1 1 auto;
    padding: 1.5rem;
    height: 100%;

    @media (min-width: size.$smMin) {
      flex: 0 0 512px;
    }

    @media (min-width: size.$lgMin) {
      flex: 0 0 480px;
    }

    &>.heading {
      max-height: 40px;
    }

    &>.container {
      height: calc(80% - 40px);

      >form {
        min-width: 100%;

        .fieldset.w-full input,
        .fieldset.w-full label {
          box-sizing: border-box;
          width: 100%;
        }
      }
    }
  }
}
</style>
