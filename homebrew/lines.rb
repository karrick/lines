require "rbconfig"
class Lines < Formula
  desc "Print a range of lines from standard input or one or more files"
  homepage "https://github.com/karrick/lines"
  version "0.5.2"

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/karrick/lines/releases/download/v0.5.2/lines_0.5.2_darwin_amd64.zip"
      sha256 "4e144cfabfb6cc564719aa8d3615abdb35e0ac8dbc0890be8c55122054db2347"
    when /linux/
      url "https://github.com/karrick/lines/releases/download/v0.5.2/lines_0.5.2_linux_amd64.tar.gz"
      sha256 "e7e069e288875f69be394cdba66104780b90c72aa8993af37609d43261745859"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  else
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/karrick/lines/releases/download/v0.5.2/lines_0.5.2_darwin_386.zip"
      sha256 "f78f84fec4e84d26a6c9b49bdf9b93bcab35b8f1b9b576ef4f5e2f6eff1fc18c"
    when /linux/
      url "https://github.com/karrick/lines/releases/download/v0.5.2/lines_0.5.2_linux_386.tar.gz"
      sha256 "78cec8ed8ae750821ff4c4a4032a097a3210cfcff6a782c9c89fe014bbd9e06e"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  end

  def install
    bin.install "lines"
  end

  test do
    system "#{bin}/lines"
  end
end
