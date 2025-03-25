import { ReactNode } from "react";

import {
  BTN_HEIGHT_LG,
  BTN_HEIGHT_MD,
  BTN_HEIGHT_SM,
  BTN_HEIGHT_XL,
  BTN_HEIGHT_XS,
} from "../../const/const";

interface ButtonProps {
  children: ReactNode;
  size?: "xs" | "sm" | "md" | "lg" | "xl";
  type?: "default" | "default-red" | "ghost" | "ghost-green" | "ghost-red";
  disabled?: boolean;
  square?: boolean;
  submit?: boolean;
  onClick?: () => void;
}

export default function Button({
  children,
  size,
  type,
  disabled = false,
  square = false,
  submit = true,
  onClick,
}: ButtonProps) {
  let btnType = "bg-primary-green text-primary-white";
  let btnHeight = BTN_HEIGHT_MD;

  if (type === "default-red") {
    btnType = "bg-primary-red text-primary-white";
  }

  if (type === "ghost") {
    btnType = "text-primary-black border border-primary-black bg-transparent";
  }

  if (type === "ghost-green") {
    btnType = "text-primary-green border border-primary-green bg-transparent";
  }

  if (type === "ghost-red") {
    btnType = "text-primary-red border border-primary-red bg-transparent";
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
      } font-semibold ${btnType} disabled:text-primary-gray disabled:bg-secondary-gray disabled:cursor-not-allowed`}
      disabled={disabled}
      type={submit ? "submit" : "button"}
      onClick={onClick}
    >
      {children}
    </button>
  );
}
