import { ReactNode } from "react";

interface ToggleProps {
  labelPosition?: "left" | "right"
  children: ReactNode;
  checked: boolean;
  setChecked: () => void;
}

export default function Toggle({ labelPosition = "left", children, checked, setChecked }: ToggleProps) {
  return (
    <label className="inline-flex items-center cursor-pointer">
      {labelPosition === "left" && <div className="mx-3">{children}</div>}
      <input
        type="checkbox"
        name=""
        value=""
        className="sr-only peer"
        checked={checked}
        onChange={setChecked}
      />
      <div className="relative w-11 h-6 bg-secondary-gray peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-secondary-green rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-primary-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-primary-white after:border-secondary-gray after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-green"></div>
      {labelPosition === "right" && <div className="mx-3">{children}</div>}
    </label>
  );
}
