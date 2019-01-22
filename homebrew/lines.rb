require 'rbconfig'
class Lines < Formula
  desc "Print a range of lines from standard input or one or more files."
  homepage "https://github.com/karrick/lines"
  version "0.5.1"

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG['host_os']
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/karrick/lines/releases/download/v0.5.1/lines_v0.5.1_darwin_amd64.zip"
      sha256 "311b3cd88f62e38a5d98ef5e5827565ebd650775db2e206b7f2d3ed2580d368e"
    when /linux/
      url "https://github.com/karrick/lines/releases/download/v0.5.1/lines_v0.5.1_linux_amd64.tar.gz"
      sha256 "b3673004e74359c42a66ba72ea2bea86aa9e6428c8ac399cd047fc4ca46fe622"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  else
    case RbConfig::CONFIG['host_os']
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/karrick/lines/releases/download/v0.5.1/lines_v0.5.1_darwin_386.zip"
      sha256 "6bbdc3986184fe277e2a06214e52cfb81db7f498a093d5b174f342abc99c2750"
    when /linux/
      url "https://github.com/karrick/lines/releases/download/v0.5.1/lines_v0.5.1_linux_386.tar.gz"
      sha256 "86baf20621821c1e3cc65ea8ec3cc6af4198be730d77880af2632e5741c7dfb3"
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
    system "lines"
  end

end
