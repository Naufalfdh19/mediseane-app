import { Facebook, Instagram, Twitter, Youtube } from "lucide-react";
import { Link } from "react-router-dom";

export default function Footer() {
  return (
    <div className="p-5 md:p-10 w-full bg-primary-green flex flex-col md:flex-row justify-around text-primary-white gap-5">
      <div className="flex flex-col gap-1">
        <p className="font-bold text-primary-white text-xl">mediSEAne</p>
        <p className="text-primary-White">
          All Your Healthcare Needs at Your Fingertips
        </p>
      </div>
      <div className="flex flex-col md:flex-row gap-5 md:gap-16">
        <div className="flex flex-col gap-1">
          <p className="font-bold text-xl">Company</p>
          <Link to="#">About</Link>
          <Link to="#">Contact</Link>
          <Link to="#">Blogs</Link>
        </div>
        <div className="flex flex-col gap-1">
          <p className="font-bold text-xl">Follow Us</p>
          <div className="flex gap-2.5 flex-wrap">
            <Facebook strokeWidth={1.7} />
            <Twitter strokeWidth={1.7} />
            <Instagram strokeWidth={1.7} />
            <Youtube strokeWidth={1.7} />
          </div>
        </div>
      </div>
    </div>
  );
}
