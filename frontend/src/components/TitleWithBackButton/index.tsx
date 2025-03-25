import { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft } from "lucide-react";

import { PATH_BACK } from "../../const/const";

interface TitleWithBackButtonProps {
  children: ReactNode;
}

export default function TitleWithBackButton({
  children,
}: TitleWithBackButtonProps) {
  const navigate = useNavigate();
  return (
    <div className="h-fit flex gap-4 items-center overflow-hidden">
      <ArrowLeft
        className="cursor-pointer flex-shrink-0"
        onClick={() => {
          navigate(PATH_BACK);
        }}
      />
      <h1 className="text-center text-2xl font-bold text-primary-black">
        {children}
      </h1>
    </div>
  );
}
