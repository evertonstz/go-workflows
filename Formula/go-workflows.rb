class GoWorkflows < Formula
    desc "Um gerenciador de workflows em Go"
    homepage "https://github.com/evertonstz/go-workflows"
    version "v0.0.6"
    license "MIT"
  
    on_macos do
      if Hardware::CPU.intel?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.6/go-workflows-darwin-amd64.tar.gz"
        sha256 "2eca6df53507bb01356a328aa2fbc370ae7760db3bc152fdc5bb773d18f41dbc"
      elsif Hardware::CPU.arm?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.6/go-workflows-darwin-arm64.tar.gz"
        sha256 "35a4a3e82736b950023506bd9fe884cf71e8d3405550a9b1af9222b098b2aab7"
      end
    end
  
    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.6/go-workflows-linux-amd64.tar.gz"
        sha256 "822c034a690a6b3eb238f90bfd72774764051f456460580f7690c2643111c092"
      elsif Hardware::CPU.arm?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.6/go-workflows-linux-arm64.tar.gz"
        sha256 "11e5b8072cc6c862351ea96995e73962b7d4b751f62b91b0359cda226b86064f"
      end
    end
  
    def install
      bin.install "go-workflows"
    end
  
    test do
      system "#{bin}/go-workflows", "--version"
    end
  end
  