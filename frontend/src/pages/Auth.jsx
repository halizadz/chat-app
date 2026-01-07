import React, { useState } from 'react';
import LoginForm from '../components/auth/LoginForm';
import RegisterForm from '../components/auth/RegisterForm';
import { MessageSquare } from 'lucide-react';

const Auth = () => {
  const [isLogin, setIsLogin] = useState(true);

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-500 via-primary-600 to-primary-700 flex items-center justify-center p-4">
      {/* Decorative background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-white/10 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-white/10 rounded-full blur-3xl" />
      </div>

      {/* Auth Card */}
      <div className="relative bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md animate-slide-up">
        {/* Logo */}
        <div className="absolute -top-12 left-1/2 -translate-x-1/2">
          <div className="w-20 h-20 bg-white rounded-2xl shadow-xl flex items-center justify-center">
            <MessageSquare className="w-10 h-10 text-primary-600" />
          </div>
        </div>

        <div className="mt-8">
          {isLogin ? (
            <LoginForm onToggle={() => setIsLogin(false)} />
          ) : (
            <RegisterForm onToggle={() => setIsLogin(true)} />
          )}
        </div>
      </div>

      {/* Footer */}
      <div className="absolute bottom-4 text-center text-white/80 text-sm">
        <p>© 2025 ChatApp. Built with ❤️ using React & Go</p>
      </div>
    </div>
  );
};

export default Auth;