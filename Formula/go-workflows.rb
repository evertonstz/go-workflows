class GoWorkflows < Formula
    desc "Um gerenciador de workflows em Go"
    homepage "https://github.com/evertonstz/go-workflows"
    version "v0.0.5"
    license "MIT"
  
    on_macos do
      if Hardware::CPU.intel?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.5/go-workflows-darwin-amd64.tar.gz"
        sha256 "00419465981e349ae338f764d49a8e0b6c847641c5da917a1c6d124bc3e983da"
      elsif Hardware::CPU.arm?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.5/go-workflows-darwin-arm64.tar.gz"
        sha256 "4155c20fa1f0f401830fea70633022f4857f806ff4dc092ab8b294c8a47c56f6"
      end
    end
  
    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.5/go-workflows-linux-amd64.tar.gz"
        sha256 "8f610346dcb2cd3421ee9c84702b3bb380221b92b4d3d9706d7152acfffbd30a"
      elsif Hardware::CPU.arm?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.5/go-workflows-linux-arm64.tar.gz"
        sha256 "d6407379335ea3b7f21613fe2b6a4b5b1ed89418ece965a58afb081a308c637e"
      end
    end
  
    def install
      bin.install "go-workflows"
    end
  
    test do
      system "#{bin}/go-workflows", "--version"
    end
  end
  