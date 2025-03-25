import { FlaskConical } from "lucide-react";
import useTitle from "../../../hooks/useTitle";
import { PHARMACIST_DASHBOARD_TITLE } from "../../../const/const";

export default function PharmacistDashboard() {
  useTitle(PHARMACIST_DASHBOARD_TITLE);

  return (
    <main className="grow flex flex-col justify-center py-24 px-16 text-center text-primary-black bg-primary-white">
      <h1 className="text-xl">Welcome</h1>
      <h2 className="text-2xl font-bold p-4">Pharmacist!</h2>
      <FlaskConical className="mx-auto my-4 w-24 h-24" strokeWidth={1.25} />
      <p className="py-6">
        Seamless management of your pharmacy products and orders
      </p>
    </main>
  );
}
