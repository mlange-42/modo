/*!
* Hooks into a *reveal.js* slide deck and allows to control SVG animations (SMIL)
* through the normal slide controls (i.e. via *reveal.js* fragments).
*
* ## Automatic fragments
*
* In your SVG, number animations with "data-fragment-index" and set begin to "indefinite":
*
* <animate data-fragment-index="1" begin="indefinite" ... />
* <animate data-fragment-index="2" begin="indefinite" ... />
*
* Fragments for controlling the animation will be inserted automatically.
*
* To properly play animations in reverse, a reverting/inverse animation can be
* added for each "normal" animation, by using negative step values:
*
* <animate data-fragment-index="1" begin="indefinite" ... />
* <animate data-fragment-index="-1" begin="indefinite" ... />
*
* To make the state after animation persistent, add `fill="freeze"` to your animations:
*
* ## Manual fragments
*
* In your SVG, mark animations with classes, and set begin to "indefinite":
*
* <animate class="first-animation" begin="indefinite" ... />
* <animate class="second-animation" begin="indefinite" ... />
*
* Include the SVG into the HTML directly, or using an <object> tag.
* Then, put empty/fake fragments on the slide with the SVG:
*
* <span class="fragment" data-svg-classes="first-animation"></span>
* <span class="fragment" data-svg-classes="second-animation"></span>
*
* A each fragment can also trigger multiple classes.
* Supports multiple SVGs per slide.
*
* Reverting/inverse animation can be added by using the class suffix "-revert":
*
* <animate class="first-animation" begin="indefinite" ... />
* <animate class="first-animation-revert" begin="indefinite" ... />
*
* ## Limitations
*
* - SVG elements to animate must have unique IDs within the document.
* - Auto-animations don't align with the ordering of other animations on the same slide.
*
* Copyright (C) 2021 Martin Lange
*/
function RevealSvgSmil() {
    var self=this;

    this.settings = {
        revert_suffix: "-revert"
    };
    this.svg_tags = [];

    Reveal.on( 'fragmentshown', event => {
        this.fragmentChanged(event.fragment, true);
    } );

    Reveal.on( 'fragmenthidden', event => {
        this.fragmentChanged(event.fragment, false);
    } );

    this.onReady = function() {
        let objs = toArray(document.getElementsByTagName("object"));
        let svgs = toArray(document.getElementsByTagName("svg"));

        self.svg_tags = objs
                    .map(obj => self.loadSvgFromObject(obj))
                    .filter(el => el != null)
                    .concat(svgs.map(obj => self.loadSvg(obj)));

        for(const entry of self.svg_tags) {
            self.initAutoAnimate(entry);
        }
    }

    this.loadSvg = function(svg) {
        let section = svg.closest("section");
        return {svg: svg, section: section, container: svg};
    }

    this.loadSvgFromObject = function(obj) {
        let doc = obj.contentDocument;
        if (doc != null) {
            let svg = doc.getElementsByTagName("svg");
            if(svg != null || svg.length > 0) {
                let section = obj.closest("section");
                return {svg: svg[0], section: section, container: obj};
            }
        }
        return null;
    }

    this.initAutoAnimate = function(svg_entry) {
        let animations = toArray(svg_entry.svg.querySelectorAll('[data-fragment-index]'));
        let anim_steps = animations.map(anim => {
            return {node: anim, step: parseInt(anim.getAttribute('data-fragment-index'))};
        });
        anim_steps.sort((a, b) => a.step - b.step);

        let steps_done = new Set();
        for(const entry of anim_steps) {
            if(entry.step >= 0) {
                let cls = "auto-animate-" + entry.step;
                entry.node.classList.add(cls);

                if(!steps_done.has(entry.step)) {
                    let fragment = self.createFragment(cls, entry.step);
                    svg_entry.container.parentElement.insertBefore(fragment, svg_entry.container);

                    steps_done.add(entry.step);
                }
            } else {
                entry.node.classList.add("auto-animate-" + Math.abs(entry.step) + self.settings.revert_suffix);
            }
        }
    }

    this.createFragment = function(cls, step) {
        let fragment = document.createElement("span");
        fragment.classList.add("fragment");
        fragment.setAttribute("data-svg-classes", cls);

        // TODO: enable after finding a solution for reveal's internal re-numbering
        // fragment.setAttribute("data-fragment-index", step);

        return fragment;
    }

    this.fragmentChanged = function(fragment, shown=true) {
        let section = fragment.closest("section");
        let index = fragment.getAttribute("data-fragment-index");

        // Unfortunately, fragment events are not fired for every single fragment, but only for one per index!
        // Thus, we need to collect all fragments with the current index.
        let fragments_now = toArray(section.getElementsByClassName("fragment")).filter(frag => {
            return index == frag.getAttribute("data-fragment-index")
        } );

        for(const frag of fragments_now) {
            if(! frag.hasAttribute('data-svg-classes')) {
                continue;
            }
            let classes = frag.getAttribute('data-svg-classes').split(/\s+/);

            for(const entry of self.svg_tags) {
                if(entry.section == section) {
                    self.toggleAnimations(entry.svg, classes, shown);
                }
            }
        }
    }

    this.toggleAnimations = function(svg, classes, begin=true) {
        for(const cls of classes) {
            let anims = svg.getElementsByClassName(cls);
            let revert = svg.getElementsByClassName(cls + self.settings.revert_suffix);

            for(const anim of anims) {
                if(begin) {
                    if(typeof anim.beginElement === 'function') {
                        anim.beginElement();
                    }
                } else {
                    if(typeof anim.endElement === 'function') {
                        anim.endElement();
                    }
                }
            }
            for(const anim of revert) {
                if(begin) {
                    if(typeof anim.beginElement === 'function') {
                        anim.endElement();
                    }
                } else {
                    if(typeof anim.beginElement === 'function') {
                        anim.beginElement();
                    }
                }
            }
        }
    }
}

function toArray(nodeList) {
    return Array.prototype.slice.call(nodeList);
}

revealSvgSmil = new RevealSvgSmil();
window.onload = revealSvgSmil.onReady;
