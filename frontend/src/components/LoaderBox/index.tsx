import Modal from "../Modal/Modal";

export default function LoaderBox() {
  return (
    <Modal>
      <div className="flex justify-center items-center gap-32">
        <div className="flex flex-row gap-2">
          <div className="w-4 h-4 rounded-full bg-primary-green animate-bounce"></div>
          <div className="w-4 h-4 rounded-full bg-primary-green animate-bounce [animation-delay:-.3s]"></div>
          <div className="w-4 h-4 rounded-full bg-primary-green animate-bounce [animation-delay:-.5s]"></div>
        </div>
      </div>
    </Modal>
  );
}
