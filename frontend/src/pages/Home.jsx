import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { MessageSquare, Users, Shield, Zap } from 'lucide-react';
import Button from '../components/ui/Button';

const Home = () => {
  const navigate = useNavigate();
  const { user } = useAuthStore();

  useEffect(() => {
    if (user) {
      navigate('/chat');
    }
  }, [user, navigate]);

  const features = [
    {
      icon: MessageSquare,
      title: 'Real-time Messaging',
      description: 'Send and receive messages instantly with WebSocket technology',
    },
    {
      icon: Users,
      title: 'Group Chats',
      description: 'Create groups and collaborate with your team seamlessly',
    },
    {
      icon: Shield,
      title: 'Secure & Private',
      description: 'Your conversations are protected with end-to-end encryption',
    },
    {
      icon: Zap,
      title: 'Fast & Reliable',
      description: 'Built with modern tech stack for lightning-fast performance',
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-500 via-primary-600 to-primary-700">
      {/* Decorative background */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-96 h-96 bg-white/10 rounded-full blur-3xl" />
        <div className="absolute top-1/2 -left-40 w-96 h-96 bg-white/10 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 right-1/3 w-96 h-96 bg-white/10 rounded-full blur-3xl" />
      </div>

      {/* Navigation */}
      <nav className="relative z-10 container mx-auto px-4 py-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <MessageSquare className="w-8 h-8 text-white" />
            <span className="text-2xl font-bold text-white">ChatApp</span>
          </div>
          <div className="flex items-center space-x-4">
            <Button
              variant="ghost"
              onClick={() => navigate('/auth')}
              className="text-white hover:bg-white/10"
            >
              Sign In
            </Button>
            <Button
              onClick={() => navigate('/auth')}
              className="bg-white text-primary-600 hover:bg-gray-100"
            >
              Get Started
            </Button>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <div className="relative z-10 container mx-auto px-4 py-20">
        <div className="text-center max-w-4xl mx-auto">
          <h1 className="text-5xl md:text-6xl font-bold text-white mb-6 animate-fade-in">
            Connect with Your Team
            <br />
            <span className="text-white/90">Anytime, Anywhere</span>
          </h1>
          <p className="text-xl text-white/90 mb-12 max-w-2xl mx-auto">
            A modern, real-time chat application built with React and Go. 
            Experience seamless communication with your team.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button
              size="lg"
              onClick={() => navigate('/auth')}
              className="bg-white text-primary-600 hover:bg-gray-100 shadow-xl"
            >
              Start Chatting Now
            </Button>
            <Button
              size="lg"
              variant="outline"
              className="border-2 border-white text-white hover:bg-white/10"
            >
              Learn More
            </Button>
          </div>
        </div>

        {/* Features Grid */}
        <div className="mt-32 grid md:grid-cols-2 lg:grid-cols-4 gap-8">
          {features.map((feature, index) => (
            <div
              key={index}
              className="bg-white/10 backdrop-blur-lg rounded-2xl p-6 text-white transform hover:scale-105 transition-transform duration-200"
            >
              <div className="w-12 h-12 bg-white/20 rounded-xl flex items-center justify-center mb-4">
                <feature.icon className="w-6 h-6" />
              </div>
              <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
              <p className="text-white/80">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Footer */}
      <div className="relative z-10 container mx-auto px-4 py-8 mt-20">
        <div className="text-center text-white/80 text-sm">
          <p>© 2025 ChatApp. Built with ❤️ using React & Go</p>
        </div>
      </div>
    </div>
  );
};

export default Home;