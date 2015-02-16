#
# Clipby Module
#
# Clipby adds a few things but not much.
# Subclass Plugin for helpers and validation methods.
#
module Clipby

  #
  # Subclass from Clipby::Plugin
  #
  class ReplacePlugin < Plugin

    def initialize
      @name = "Replace Plugin Plugin"

      @description <<-description
      replace e with word
      description

      @author = "Dustin Willis Webber"
      @version = "0.1.0"
    end

    def self.some_network_stuff
      s = TCPSocket.open("intel.criticalstack.com", 80)
      s.write("GET / HTTP/1.0\r\n\r\n")
      log s.read
      s.close
    end

    def self.run(type, data)
      decode = Base64.decode(data)
      return type, decode.gsub!("e", "(づ｡◕‿‿◕｡)づ ︵ ┻━┻")
    end

  end

end
