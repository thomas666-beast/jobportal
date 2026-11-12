export function generatePassword(length = 16) {
  const charset =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
  const array = new Uint32Array(length)
  crypto.getRandomValues(array)

  return Array.from(array, (n) => charset[n % charset.length]).join('')
}
