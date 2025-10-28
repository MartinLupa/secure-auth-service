import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { NextRequest, NextResponse } from "next/server"

export async function GET(req: NextRequest) {
 const cookieStore = await cookies()
 const authToken = req.cookies.get("redirect_session_token")

 if (!authToken) {
  return NextResponse.json({ error: "No auth token found" }, { status: 401 })
 }

 cookieStore.set("session_token", authToken.value)
 cookieStore.delete("redirect_session_token")

 redirect("/logged")
}