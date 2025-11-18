import LoginForm from "@/components/ui/auth/login-form";

export default function LoginPage() {
  return (
    <div className="relative bg-background w-full h-screen overflow-hidden">
      <div
        className="absolute inset-0 bg-red-500 bg-cover"
        style={{
          backgroundImage: "url('/login-unsch.jpg')",
        }}
      >
        <div className="bg-black/20 w-full h-full"></div>
      </div>

      <div
        className="absolute inset-0 bg-muted"
        style={{
          clipPath: "polygon(70% 0, 100% 0, 100% 100%, 65% 100%)",
        }}
      ></div>

      <div className="z-10 relative flex justify-center items-center p-4 w-full h-full">
        <div className="relative rounded-[60px] w-6xl h-4/6 overflow-hidden">
          <div
            className="absolute inset-0 bg-cover bg-center"
            style={{
              backgroundImage: "url('/login-unsch.jpg')",
            }}
          ></div>
          <div
            className="absolute inset-0 bg-background z-20 flex justify-end items-center"
            style={{
              clipPath: "polygon(60% 0, 100% 0, 100% 100%, 55% 100%)",
            }}
          >
            <div className="w-[40%] flex justify-center items-center">
              <LoginForm />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
