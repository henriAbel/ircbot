/*
 * Very lightweight jquery lightbox plugin. Supports custom content eq youtube embeded videos.
 * 
 * by Henri Abel
 *
 * More information in:
 * http://www.github.com/henriAbel/
 */

(function ($) {

    $.fn.popup = function(options) {

        var settings = $.extend({
            getContentUrl: function(e) {
                return e.attr('data-url');
            }
        }, options );

        $this = $(this);
        var container,
            content;
        var opened = false;
        var active;

        // Remove old bindings if exists
        this.unbind('click.popup');

        var showContainer = function() {
            if (!opened) {
                beforeOpen();
                container.show();
                opened = true;    
            }
        }

        var closeContainer = function() {
            container.hide();
            content.empty();
            opened = false;
            afterClose();
        }

        var hasNext = function() {
            return getNext().length > 0;
        }

        var getNext = function() {
            return active.next($this.selector);
        }

        var hasPrevious = function() {
            return getPrevious().length > 0;
        }

        var getPrevious = function() {
            return active.prev($this.selector);
        }

        var setActive = function(element) {
            active = element;
            url = settings.getContentUrl(active);
            content.html('<img src="' + url + '"/>');
            showContainer();
        }

        var createContainer = function() {
            var old = document.getElementById('popup-container');
            if (null != old) {
                old.parentNode.removeChild(old);    
            }

            econtainer = document.createElement('div');
            econtainer.setAttribute('class', 'text-center');
            econtainer.setAttribute('id', 'popup-container');
            estyle = econtainer.style;
            estyle.display = 'none';
            estyle.width = '100%';
            estyle.height = '100%';
            estyle.position = 'absolute';
            estyle.top = 0;

            c = document.createElement('div');
            style = c.style;
            style.display = 'inline-block';
            style.zIndex = '10001';
            style.position = 'relative';
            style.left = '0';
            style.right = '0';
         
            econtainer.appendChild(c);

            b = document.createElement('div');
            b.setAttribute('id', 'popup-background');
            bstyle = b.style;
            bstyle.backgroundColor = 'black';
            bstyle.zIndex = '10000'; // Twitter bootstrap menus are z-index 1000
            bstyle.opacity = '.5';
            bstyle.width = '100%';
            bstyle.height = '100%';
            bstyle.top = 0;
            bstyle.position = 'absolute';
            bstyle.filter = 'alpha(opacity=50)'; //For IE

            econtainer.appendChild(b);
            document.body.appendChild(econtainer);

            container = $(econtainer);
            content = $(c);
        }

        var beforeOpen = function() {
            $(document).bind('keyup.popup', function(e) {
                switch(e.which) {
                    case 27: //Esc
                        closeContainer();
                        break;
                    case 39: // Arrow right
                        if (hasNext()) {
                            setActive(getNext());    
                        }
                        break;
                    case 37: // Arrow left
                        if (hasPrevious()) {
                            setActive(getPrevious());    
                        }
                        break;
                    }
            });

            $(document).bind('click.popupOutside', function(e) {
                if (e.target.getAttribute('id') == 'popup-background') {
                    closeContainer();
                }
            });
        }

        var afterClose = function() {
            $(document).unbind('keyup.popup');
            $(document).unbind('click.popupOutside');
        }
        
        createContainer();
         
        this.bind('click.popup', function(e) {
            setActive($(e.target));
        });

    };
}(jQuery));