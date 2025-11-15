import UserLogger from "@/components/UserLogger";
import Image from "next/image";

export default function Home() {
  return (
   <main>
    <h1>Dashboard</h1>
    <UserLogger />
   </main>
  );
}
