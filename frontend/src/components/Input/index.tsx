import { ReactNode } from "react";

interface inputProps {
  type: "text" | "password" | "number" | "time";
  name: string;
  placeholder: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  icon?: ReactNode;
  onIconClick?: () => void;
  valid?: boolean;
  value?: string;
  round?: boolean;
  readOnly?: boolean;
}

export default function Input({
  type,
  name,
  placeholder,
  onChange,
  icon,
  onIconClick,
  valid = true,
  value,
  round = false,
  readOnly,
}: inputProps) {
  return (
    <div
      className={`w-full flex border px-[13px] py-[12px] ${
        round ? "rounded-full" : "rounded-lg"
      } gap-2 
      ${readOnly ? "bg-secondary-gray cursor-not-allowed" : "bg-primary-white"}
      ${valid ? "border-primary-black" : "border-primary-red"}`}
    >
      <input
        aria-label="input"
        type={type}
        name={name}
        id={name}
        value={value}
        placeholder={placeholder}
        onChange={onChange}
        readOnly={readOnly}
        className={`min-w-0 grow focus:outline-none text-primary-black ${
          readOnly ? "bg-secondary-gray cursor-not-allowed" : "bg-primary-white"
        }`}
      />
      {icon && (
        <div
          onClick={onIconClick}
          className="cursor-pointer text-primary-black"
        >
          {icon}
        </div>
      )}
    </div>
  );
}
