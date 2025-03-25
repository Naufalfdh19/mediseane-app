import { Send } from "lucide-react";

export default function ResetLinkSent() {
  return (
    <div className="grow flex bg-primary-white">
      <div className="w-[90%] lg:w-[80%] mx-auto flex flex-col justify-center items-center gap-3">
        <Send
          className="h-32 w-32 md:h-64 md:w-64 text-primary-green"
          strokeWidth={0.5}
        />
        <h1 className="text-center text-2xl md:text-4xl font-bold text-primary-black">
          Password Reset Link Sent
        </h1>
        <p className="text-center md:text-xl text-primary-black">
          A password reset link has been sent to your email address
        </p>
      </div>
    </div>
  );
}
