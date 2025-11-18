/* eslint-disable @next/next/no-img-element */
"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Eye, EyeOff } from "lucide-react";

export default function LoginForm() {
  const [showPassword, setShowPassword] = useState(false);

  return (
    <div className="flex justify-center items-center w-full max-w-md p-6">
      <div className="w-full">
        <h1 className="mb-2 font-bold text-gray-900 text-4xl text-center">
          Bienvenido
        </h1>
        <p className="mb-8 text-gray-500 text-center">
          Inicia sesión para continuar
        </p>

        <div className="flex flex-col gap-4 mb-4">
          <Input
            type="email"
            placeholder="Correo electrónico"
            className="rounded-xl h-12 text-base"
          />

          <div className="relative">
            <Input
              type={showPassword ? "text" : "password"}
              placeholder="Contraseña"
              className="pr-12 rounded-xl h-12 text-base"
            />

            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="top-1/2 right-4 absolute text-gray-500 hover:text-gray-700 -translate-y-1/2"
            >
              {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
            </button>
          </div>
        </div>

        <div className="flex justify-end mb-4">
          <button className="text-red-500 text-sm hover:underline">
            ¿Olvidaste tu contraseña?
          </button>
        </div>

        <div className="flex items-center gap-4 my-6">
          <Separator className="flex-1" />
          <span className="text-gray-500 text-sm">o</span>
          <Separator className="flex-1" />
        </div>

        <Button className="bg-red-500 hover:bg-red-600 rounded-xl w-full h-12 text-lg">
          Iniciar sesión
        </Button>

        <Button
          variant="outline"
          className="flex items-center gap-2 mt-4 rounded-xl w-full h-12"
        >
          <img src="/icons/google-icon.svg" alt="Google" className="w-5 h-5" />
          Iniciar con Google
        </Button>

        <p className="mt-6 text-gray-500 text-sm text-center">
          ¿No tienes una cuenta?{" "}
          <span className="font-semibold text-red-500 hover:underline cursor-pointer">
            Regístrate
          </span>
        </p>
      </div>
    </div>
  );
}
