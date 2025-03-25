import { CircleAlert } from "lucide-react";

import Modal from "../Modal/Modal";
import Button from "../Button";
import { ReactNode } from "react";

interface ErrorBoxProps {
  children: ReactNode;
  onClose?: () => void;
}

export default function ErrorBox({
  children,
  onClose = () => {},
}: ErrorBoxProps) {
  return (
    <Modal onClose={onClose}>
      <div className="w-[80%] mx-auto flex flex-col items-center justify-center gap-6">
        <CircleAlert className="w-[100px] h-[100px] md:w-[130px] md:h-[130px] text-primary-red" />
        <div className="text-center flex flex-col gap-3">{children}</div>
        <div className="w-[130px] md:w-[163px]">
          <Button submit={false} type="ghost" size="md" onClick={onClose}>
            Close
          </Button>
        </div>
      </div>
    </Modal>
  );
}
