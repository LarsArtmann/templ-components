// Lightbox component: a native <dialog> image viewer with thumbnail
// triggers, prev/next navigation, rotate and zoom — extracted from
// DiscordSync's image thumbnails and nsfw-classifier's image viewer
// (2026-10-10 analysis, TODO_LIST #396). Zero framework JS: one CSP-safe
// singleton script; Escape, focus trap, and backdrop dismissal come from the
// native dialog.
package display

import (
	"github.com/larsartmann/templ-components/utils"
)

// LightboxImage is one viewable image.
type LightboxImage struct {
	// Src is the full-size image URL.
	Src string
	// Alt is the image's alternative text (required for a11y).
	Alt string
	// Caption optionally renders under the image.
	Caption string
	// ThumbSrc optionally overrides the thumbnail source (TriggerLabel
	// button mode uses Src when empty).
	ThumbSrc string
}

// LightboxProps configures the lightbox.
type LightboxProps struct {
	utils.BaseProps

	// Images renders in order; the viewer navigates them with wrap-around.
	Images []LightboxImage
	// TriggerLabel renders one opener button with this text (Thumbnails
	// overrides it). Empty + no Thumbnails renders the dialog only —
	// open it yourself via tcOpenLightbox("<id>-dialog", index).
	TriggerLabel string
	// Thumbnails renders one opener button per image (using ThumbSrc/Src
	// as its thumbnail) instead of a single TriggerLabel button.
	Thumbnails bool
}

// DefaultLightboxProps returns sensible defaults.
func DefaultLightboxProps() LightboxProps {
	return LightboxProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		TriggerLabel: "View images",
	}
}

// lightboxJS is the CSP-safe singleton script: open/navigate/rotate/zoom on
// [data-tc-lightbox] roots via event delegation. Escape/backdrop dismissal
// are native dialog behavior; rotate and zoom compose on the image's
// transform (JS-set styles are not blocked by style-src).
const lightboxJS = `if(!window.tcLightboxAttached){window.tcLightboxAttached=true;` +
	`window.tcLightboxApply=function(img){` +
	`var r=img.getAttribute('data-tc-lightbox-rot')||'0';` +
	`var z=img.getAttribute('data-tc-lightbox-zoom')==='1'?' scale(2)':'';` +
	`img.style.transform='rotate('+r+'deg)'+z;};` +
	`window.tcShowLightbox=function(dlg,idx){` +
	`var imgs=dlg.querySelectorAll('[data-tc-lightbox-img]');if(!imgs.length)return;` +
	`if(idx<0)idx=imgs.length-1;if(idx>=imgs.length)idx=0;` +
	`for(var k=0;k<imgs.length;k++){imgs[k].hidden=k!==idx;}` +
	`var cap=dlg.querySelector('[data-tc-lightbox-caption]');` +
	`if(cap)cap.textContent=imgs[idx].getAttribute('data-tc-lightbox-cap')||'';` +
	`var img=imgs[idx];img.style.transform='';img.setAttribute('data-tc-lightbox-rot','0');img.setAttribute('data-tc-lightbox-zoom','0');img.classList.remove('cursor-zoom-out');` +
	`dlg.setAttribute('data-tc-lightbox-index',String(idx));` +
	`if(!dlg.open)dlg.showModal();};` +
	`window.tcOpenLightbox=function(id,idx){var d=document.getElementById(id);if(d)window.tcShowLightbox(d,idx||0);};` +
	`document.addEventListener('click',function(e){` +
	`var opener=e.target.closest('[data-tc-lightbox-open]');` +
	`if(opener){var root=opener.closest('[data-tc-lightbox]');if(!root)return;` +
	`var dlg=document.getElementById(root.id+'-dialog');if(!dlg)return;` +
	`window.tcShowLightbox(dlg,parseInt(opener.getAttribute('data-tc-lightbox-open'),10)||0);return;}` +
	`var nav=e.target.closest('[data-tc-lightbox-prev],[data-tc-lightbox-next]');` +
	`if(nav){var dr=nav.closest('dialog');if(!dr)return;` +
	`var cur=parseInt(dr.getAttribute('data-tc-lightbox-index')||'0',10);` +
	`window.tcShowLightbox(dr,cur+(nav.hasAttribute('data-tc-lightbox-next')?1:-1));return;}` +
	`var rot=e.target.closest('[data-tc-lightbox-rotate]');` +
	`if(rot){var rr=rot.closest('dialog');if(!rr)return;` +
	`var ri=rr.querySelectorAll('[data-tc-lightbox-img]')[parseInt(rr.getAttribute('data-tc-lightbox-index')||'0',10)];` +
	`if(ri){var r=(parseInt(ri.getAttribute('data-tc-lightbox-rot')||'0',10)+90)%360;` +
	`ri.setAttribute('data-tc-lightbox-rot',String(r));tcLightboxApply(ri);}return;}` +
	`var zoom=e.target.closest('[data-tc-lightbox-zoom]');` +
	`if(zoom){var zr=zoom.closest('dialog');if(!zr)return;` +
	`var zi=zr.querySelectorAll('[data-tc-lightbox-img]')[parseInt(zr.getAttribute('data-tc-lightbox-index')||'0',10)];` +
	`if(zi){zi.setAttribute('data-tc-lightbox-zoom',zi.getAttribute('data-tc-lightbox-zoom')==='1'?'0':'1');` +
	`zi.classList.toggle('cursor-zoom-out');tcLightboxApply(zi);}}` +
	`});}` +
	`(function(){var dlgs=document.querySelectorAll('dialog[data-tc-lightbox-dialog]');` +
	`for(var i=0;i<dlgs.length;i++){dlgs[i].addEventListener('click',function(e){` +
	`if(e.target.closest('[data-tc-close]')||e.target===this)this.close();});}})();`
