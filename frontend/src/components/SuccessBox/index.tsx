import { ReactNode } from "react";
import { CircleCheck } from "lucide-react";

import Modal from "../Modal/Modal";
import Button from "../Button";

interface SuccessBoxProps {
  children: ReactNode;
  onClose?: () => void;
}

export default function SuccessBox({
  children,
  onClose = () => {},
}: SuccessBoxProps) {
  return (
    <Modal onClose={onClose}>
      <div className="w-full py-8 px-2 flex flex-col items-center justify-center gap-6">
        <CircleCheck className="w-[100px] h-[100px] md:w-[130px] md:h-[130px] text-primary-green" />
        <h1 className="w-[95%] font-semibold text-primary-black text-2xl text-center">
          {children}
        </h1>
        <div className="w-[130px] md:w-[163px]">
          <Button submit={false} type="ghost" size="md" onClick={onClose}>
            Close
          </Button>
        </div>
      </div>
    </Modal>
  );
}
