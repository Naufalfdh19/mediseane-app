import { ReactNode } from "react";

interface textareaProps {
  name: string;
  placeholder: string;
  onChange?: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
  icon?: ReactNode;
  onIconClick?: () => void;
  valid?: boolean;
  value?: string;
}

export default function TextArea({
  name,
  placeholder,
  onChange,
  icon,
  onIconClick,
  valid = true,
  value,
}: textareaProps) {
  return (
    <div
      className={`w-full flex border px-[13px] py-[12px] rounded-lg gap-2 ${
        valid ? "border-primary-black" : "border-primary-red"
      }`}
    >
      <textarea
        name={name}
        id={name}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        className="min-w-0 grow focus:outline-none text-primary-black bg-primary-white"
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
