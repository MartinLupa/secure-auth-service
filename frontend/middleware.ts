import { NextRequest, NextResponse } from "next/server";

const protectedPaths = ['/logged']

export async function middleware(request: NextRequest) {
 if (protectedPaths.some((path) => request.nextUrl.pathname.startsWith(path))) {
  const authCookie = request.cookies.get('session_token')

  if (!authCookie) {
   const loginUrl = new URL('/login', request.url)
   return Response.redirect(loginUrl.toString())
  }

  try {
   console.log("authCookie: ", authCookie)
   const response = await fetch('http://localhost:8000/jwt/validate', {
    method: 'POST',
    headers: {
     'Content-Type': 'application/json',
     'Authorization': `${authCookie.value}`
    }
   })

   if (!response.ok) {
    const loginUrl = new URL('/login', request.url)
    return Response.redirect(loginUrl.toString())
   }

  } catch (error) {
   console.error('Error validating JWT:', error)
   const loginUrl = new URL('/login', request.url)
   return Response.redirect(loginUrl.toString())
  }

  return NextResponse.next()
 }
}