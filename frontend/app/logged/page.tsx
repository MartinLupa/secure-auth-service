import { FaCheck } from 'react-icons/fa';

export default function Page() {
 return (
  <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
   <div className="w-full max-w-sm text-center">
    <FaCheck className="mx-auto mb-4 text-4xl text-green-500" />
    <h1 className="mb-2 text-2xl font-bold">Log in Successful!</h1>
   </div>
  </div>
 )
}