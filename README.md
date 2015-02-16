# Clipby - clipboard manipulation and organization

Clipby will monitor your clipboard and organize copies by type. The clipboard data is stored in an embeded db for 
easy access. Clipby also embeds mruby for custom scripts.

# Ruby Script Example

This plugin will replace all copied data that contains the letter `e` with `ZZ`. Obviosuly this is a stupid example
but you could do cool stuff like .. copy json, upload to gist.. return a new type and data to replace the contents of
the clipboard.

The run function must return a type and data.

```ruby
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
```

![plugin output](https://raw.githubusercontent.com/mephux/clipby/master/plugin-exaple-output.png)

# Build

  * make sure you have godep installed.
  * clone the go-mruby repo.
  * build mruby using the build_config.rb in the clipby root dir.
  * make sure to add gems you make need - by default i add base64, socket and pack.
  * copy the libmruby.a to the root of the clipby dir.
  * make
  * Enjoy!



# Notes
  
  * THIS IS BETA
  * This is Mac OSX only.
  * Clipby will create $HOME/.clipby & $HOME/.clipby/plugins/
  * Any .rb file in plugins will be loaded automatically.
