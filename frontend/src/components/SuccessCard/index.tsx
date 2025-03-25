interface SuccessCardProps {
  message: string;
}

export default function SuccessCard({ message }: SuccessCardProps) {
  return (
    <div className="w-full py-3 text-center text-primary-green bg-secondary-green border rounded-lg">
      {message}
    </div>
  );
}
