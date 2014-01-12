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
            },
            fittocontent: true
        }, options );

        $this = $(this);
        var container,
            content,
            dummy,
            wrapper;
        var opened = false;
        var active;

        // Remove old bindings if exists
        this.unbind('click.popup');

        var showContainer = function() {
            if (!opened) {
                beforeOpen();
                container.show();
                container.css('opacity', '1');
                opened = true;    
            }
        }

        var closeContainer = function() {
            container.css('opacity', '0');
            content.empty();
            setTimeout(function() {
                container.css('display', 'none');
                wrapper.css({
                    'width':     '160px',
                    'height':    '100px'
                });
            }, 300);
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
            if (opened) {
                // TODO add some loading gif or smth
                content.empty();
                dummy.show();
            }
            preload(url, function(e) {
                e.style.maxWidth = '100%';
                e.style.maxHeight = '100%';
                // Check if need to scale wrapper 
                var w = $(window).width() * 0.9;
                var h = $(window).height() * 0.9;
                var pw = this.width;
                var ph = this.height;
                var ratio = pw / ph;

                if (ph > h) {
                    console.log('windowW: ' + w);
                    console.log('windowH: ' + h);
                    console.log('pw: ' + pw);
                    console.log('ph: ' + ph);
                    console.log('ratio: ' + ratio);
                    ph = h;
                    pw = ratio * h;
                }
                
                if (pw > w) {
                    pw = w;
                    ph = pw / ratio;
                }

                wrapper.css({
                    'width':     pw,
                    'height':    ph
                });

                setTimeout(function() {
                    content.html(e);
                    dummy.hide();
                }, 300);
            });
            
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
            estyle.zIndex = '10000'; // Twitter bootstrap menus are z-index 1000
            //estyle.display = 'none';
            estyle.width = '100%';
            estyle.height = '100%';
            estyle.position = 'fixed';
            estyle.top = 0;
            estyle.opacity = '0';
            estyle.display = 'none';
            estyle.filter = 'alpha(opacity=0)'; //For IE
            doVendor(estyle, 'all 0.3s ease');

            c = document.createElement('div');
            style = c.style;
            style.display = 'inline-block';
            style.zIndex = '10001';
            style.position = 'absolute';
            style.left = '0';
            style.top = '0';
            style.right = '0';
            style.bottom = '0';
            style.margin = 'auto';
            style.maxWidth = '90%';
            style.maxHeight = '90%';
            style.width = '160px';
            style.height = '100px';
            doVendor(style, 'all 0.3s ease');

            dummy = document.createElement('div');
            dummy.setAttribute('id', 'dummy');
            dstyle = dummy.style;
            dstyle.backgroundColor = '#ccc';
            dstyle.backgroundImage = "url('data:image/gif;base64,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')";
            dstyle.backgroundPosition = '50% 50%';
            dstyle.backgroundRepeat = 'no-repeat';
            dstyle.position = 'absolute';
            dstyle.zIndex = '-1';
            dstyle.width = '100%';
            dstyle.height = '100%';

            ccontainer = document.createElement('div');

            econtainer.appendChild(c);
            c.appendChild(dummy);
            c.appendChild(ccontainer);

            b = document.createElement('div');
            b.setAttribute('id', 'popup-background');
            bstyle = b.style;
            bstyle.backgroundColor = 'black';
            bstyle.opacity = '.7';
            bstyle.width = '100%';
            bstyle.height = '100%';
            bstyle.top = 0;
            bstyle.position = 'absolute';
            bstyle.filter = 'alpha(opacity=70)'; //For IE

            econtainer.appendChild(b);
            document.body.appendChild(econtainer);

            container = $(econtainer);
            wrapper = $(c);
            content = $(ccontainer);
            dummy = $(dummy);
        }

        var doVendor = function(element, content) {
            element.webkitTransition = content;
            element.MozTransition = content;
            element.msTransition = content;
            element.OTransition = content;
            element.transition = content;
        }

        var preload = function(src, callback) {
            image = new Image();
            image.onload = function() {
                callback.call(this, image);
            }
            image.onerror = function() {
                callback.call(this, image);
            }
            image.src = src;
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