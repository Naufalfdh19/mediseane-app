import { ReactNode } from "react";

import {
  BTN_HEIGHT_LG,
  BTN_HEIGHT_MD,
  BTN_HEIGHT_SM,
  BTN_HEIGHT_XL,
  BTN_HEIGHT_XS,
} from "../../const/const";

interface CategoryButtonProps {
  children: ReactNode;
  size?: "xs" | "sm" | "md" | "lg" | "xl";
  type?: "category" | "classification";
  disabled?: boolean;
  square?: boolean;
  submit?: boolean;
  onClick?: () => void;
}

export default function CategoryButton({
  children,
  size,
  type = "category",
  disabled = false,
  square = false,
  submit = true,
  onClick,
}: CategoryButtonProps) {
  let btnType = "text-primary-green bg-primary-green bg-opacity-20";
  let btnHeight = BTN_HEIGHT_MD;

  if (type === "classification") {
    btnType = "text-primary-blue bg-primary-blue bg-opacity-20";
  }

  if (size === "xs") {
    btnHeight = BTN_HEIGHT_XS;
  }
  if (size === "sm") {
    btnHeight = BTN_HEIGHT_SM;
  }
  if (size === "lg") {
    btnHeight = BTN_HEIGHT_LG;
  }
  if (size === "xl") {
    btnHeight = BTN_HEIGHT_XL;
  }

  return (
    <button
      className={`w-full ${btnHeight} ${
        square ? "rounded" : "rounded-full"
      } font-semibold ${btnType} disabled:text-primary-gray disabled:bg-secondary-gray`}
      disabled={disabled}
      type={submit ? "submit" : "button"}
      onClick={onClick}
    >
      {children}
    </button>
  );
}
