class GoWorkflows < Formula
    desc "Um gerenciador de workflows em Go"
    homepage "https://github.com/evertonstz/go-workflows"
    version "v0.0.7"
    license "MIT"
  
    on_macos do
      if Hardware::CPU.intel?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.7/go-workflows-darwin-amd64.tar.gz"
        sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
      elsif Hardware::CPU.arm?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.7/go-workflows-darwin-arm64.tar.gz"
        sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
      end
    end
  
    on_linux do
      if Hardware::CPU.intel?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.7/go-workflows-linux-amd64.tar.gz"
        sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
      elsif Hardware::CPU.arm?
        url "https://github.com/evertonstz/go-workflows/releases/download/v0.0.7/go-workflows-linux-arm64.tar.gz"
        sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
      end
    end
  
    def install
      bin.install "go-workflows"
    end
  
    test do
      system "#{bin}/go-workflows", "--version"
    end
  end
  