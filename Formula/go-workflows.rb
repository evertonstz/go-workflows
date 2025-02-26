class GoWorkflows < Formula
    desc "Um gerenciador de workflows em Go"
    homepage "https://github.com/evertonstz/go-workflows"
    version "0.0.5"
    license "MIT"
  
    on_macos do
      if Hardware::CPU.intel?
        url "https://github.com/evertonstz/go-workflows/releases/download/0.0.5/go-workflows-darwin-amd64.tar.gz"
        sha256 "INSIRA_SHA256_DARWIN_AMD64"
      elsif Hardware::CPU.arm?
        url "https://github.com/evertonstz/go-workflows/releases/download/0.0.5/go-workflows-darwin-arm64.tar.gz"
        sha256 "INSIRA_SHA256_DARWIN_ARM64"
      end
    end
  
    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/evertonstz/go-workflows/releases/download/0.0.5/go-workflows-linux-amd64.tar.gz"
        sha256 "INSIRA_SHA256_LINUX_AMD64"
      elsif Hardware::CPU.arm?
        url "https://github.com/evertonstz/go-workflows/releases/download/0.0.5/go-workflows-linux-arm64.tar.gz"
        sha256 "INSIRA_SHA256_LINUX_ARM64"
      end
    end
  
    def install
      bin.install "go-workflows"
    end
  
    test do
      system "#{bin}/go-workflows", "--version"
    end
  end
  