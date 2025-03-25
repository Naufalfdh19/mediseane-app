import { ReactNode } from "react";

interface AuthFormContainerProps {
  children: ReactNode;
}

export default function AuthFormContainer({
  children,
}: AuthFormContainerProps) {
  return (
    <div className="bg-primary-white grow flex justify-center items-center">
      <div className="w-[90%] lg:w-[80%] mx-auto flex justify-center items-center gap-32">
        <div className="hidden lg:block lg:w-[50%] text-center ">
          <p className="lg:text-6xl xl:text-8xl font-bold text-primary-black">
            medi<span className="text-primary-green">SEA</span>ne
          </p>
          <p className="text-primary-black">
            All Your <span className="text-primary-green">Healthcare</span>{" "}
            Needs at Your Fingertips
          </p>
        </div>
        <div className="w-full  max-w-[400px] flex flex-col gap-[22px]">
          {children}
        </div>
      </div>
    </div>
  );
}
