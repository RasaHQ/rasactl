class Rasactl < Formula
  desc "rasactl"
  homepage "https://github.com/RasaHQ/rasactl"
  version "0.0.4"

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/RasaHQ/rasactl/releases/download/0.0.4/rasactl_0.0.4_darwin_amd64.tar.gz"
    sha256 "960505b7992669fb8d6dadcc77e6f2a0f4ca0f4e78c6fc674b332aba9522e845"
  end
  if OS.mac? && && Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
    url "https://github.com/RasaHQ/rasactl/releases/download/0.0.4/rasactl_0.0.4_darwin_arm64.tar.gz"
    sha256 "5c51efc7e802b4fc4635b4c3305c7aebe310eb5ba6653188626098a8a0d1b288"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/RasaHQ/rasactl/releases/download/0.0.4/rasactl_0.0.4_linux_amd64.tar.gz"
    sha256 "358c8badd76895e189a03e4bdc326e71abda2697b65c06c659521d36e338db1b"
  end

  def install
    bin.install "rasactl"
  end

end
