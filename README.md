# Clipby - clipboard manipulation and organization

Clipby will monitor your clipboard and organize copies by type. The clipboard data is stored in an embeded db for 
easy access. Clipby also embeds mruby for custom scripts.

# Ruby Script Example

This plugin will replace all copied data that contains the letter `e` with `ZZ`. Obviosuly this is a stupid example
but you could do cool stuff like .. copy json, upload to gist.. return a new type and data to replace the contents of
the clipboard.

The run function must return a type and data.

```ruby
module Clipby

  class ReplacePlugin < Plugin

    def initialize
      @name = "Replace Plugin"

      @description <<-description
      Replace the letter e with ZZ from all data put into the system clipboard.
      description

      @author = "Dustin Willis Webber"
      @version = "0.1.0"
    end

    def self.run(type, data)
      log "Got Data!"
      return type, data.gsub!("e", "ZZ")
    end

  end

end
```

![plugin output]()

# Notes
  
  * THIS IS BETA
  * This is Mac OSX only.
  * Clipby will create $HOME/.clipby & $HOME/.clipby/plugins/
  * Any .rb file in plugins will be loaded automatically.
