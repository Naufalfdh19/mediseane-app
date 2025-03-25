import { ReactNode } from "react";
interface ModalProps {
  onClose?: () => void;
  children?: ReactNode;
}

export default function Modal({ onClose = () => {}, children }: ModalProps) {
  return (
    <>
      <div
        className="fixed flex justify-center items-center w-screen h-screen bg-primary-gray bg-opacity-50 top-0 left-0"
        onClick={() => {
          onClose();
        }}
      >
        <div className="w-[80%] min-h-[364px] mx-auto py-6 md:w-[631px] md:min-h-[450px] flex justify-center items-center bg-primary-white rounded-lg" onClick={e => e.stopPropagation()}>
          {children}
        </div>
      </div>
    </>
  );
}
