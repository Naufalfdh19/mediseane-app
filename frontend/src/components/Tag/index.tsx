import { ReactNode } from "react";

interface TagProps {
  children: ReactNode;
}

export default function Tag({ children }: TagProps) {
  return (
    <div className="w-fit bg-secondary-green px-5 py-2 rounded-full">
      <p className="text-primary-green font-semibold text-lg">{children}</p>
    </div>
  );
}
