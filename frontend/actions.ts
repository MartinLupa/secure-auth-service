import { redirect } from "next/navigation"

export async function loginAction(prevState: any, formData: FormData) {
 const email = formData.get('email') as string
 const password = formData.get('password') as string

 if (!email || !password) {
  return { error: 'Email and password are required.' }
 }

 try {
  const response = await fetch(`http://localhost:8000/login`, {
   method: 'POST',
   headers: { 'Content-Type': 'application/json' },
   body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
   const error = await response.json()
   return { error: error.error || 'Login failed' }
  }

  const data = await response.json()

  return { success: true, data }
 } catch (error) {
  return { error: 'Network error. Please try again.' }
 }
}

export async function signupAction(prevState: any, formData: FormData) {
 const fullName = formData.get('full-name') as string
 const email = formData.get('email') as string
 const password = formData.get('password') as string
 const confirmPassword = formData.get('confirm-password') as string

 if (!fullName || !email || !password || !confirmPassword) {
  return { error: 'All fields are required.' }
 }

 if (password !== confirmPassword) {
  return { error: 'Passwords do not match.' }
 }

 try {
  const response = await fetch(`http://localhost:8000/signup`, {
   method: 'POST',
   headers: { 'Content-Type': 'application/json' },
   body: JSON.stringify({ full_name: fullName, email, password, confirm_password: confirmPassword }),
  })

  if (!response.ok) {
   const error = await response.json()
   return { error: error.error || 'Signup failed' }
  }

  return { success: true }

 } catch (error) {
  return { error: 'Network error. Please try again.' }
 }
}

export async function verifyOTPAction(prevState: any, formData: FormData) {
 const otpCode = formData.get('otp') as string

 if (!otpCode) {
  return { error: 'OTP code is required.' }
 }

 console.log("OTP code from action: ", otpCode)

 try {
  const response = await fetch(`http://localhost:8000/otp/verify`, {
   method: 'POST',
   headers: { 'Content-Type': 'application/json' },
   body: JSON.stringify({ otp_code: otpCode }),
  })

  if (!response.ok) {
   const error = await response.json()
   return { error: error.error || 'OTP verification failed' }
  }

  return { success: true }

 } catch (error) {
  return { error: 'It was not possible to verify the OTP code. Please try again or click on resend code.' }
 }
}