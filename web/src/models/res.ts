export class Res<T> {
  status: number = 200
  message: string = ""
  data: T | null = null
}