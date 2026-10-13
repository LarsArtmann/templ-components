package display

// carouselJS is the carousel's singleton script: click delegation for the
// arrow buttons and dots, keyboard navigation (RTL-aware), and scrollend
// (with a scroll+debounce fallback) dot sync. Extracted verbatim from the
// former inline <script> block; rendered via utils.ScriptComponent so an
// empty nonce OMITS the attribute (never nonce="") without dropping the
// script for non-CSP consumers — the omit-empty rule (issues #7/#9/#24).
const carouselJS = `			(function(){
			if(window.tcCarouselAttached)return;window.tcCarouselAttached=true;
			var activeDot='h-2 w-2 rounded-full bg-blue-600 dark:bg-blue-500';
			var idleDot='h-2 w-2 rounded-full bg-gray-400/60 dark:bg-gray-500/60';
			function tcCarouselTrack(c){return c.querySelector('[data-tc-carousel-track]');}
			function tcCarouselIndex(track){return Math.round(track.scrollLeft/track.offsetWidth);}
			function tcCarouselCount(track){return Math.max(1,Math.round(track.scrollWidth/track.offsetWidth));}
			function tcUpdateDots(track){
			var c=track.closest('[data-tc-carousel]');if(!c)return;
			var idx=tcCarouselIndex(track);
			var dots=c.querySelectorAll('[data-tc-carousel-dot]');
			for(var i=0;i<dots.length;i++){var d=dots[i].firstElementChild||dots[i];d.className=i===idx?activeDot:idleDot;}
			}
			function tcCarouselGo(c,idx){
			var track=tcCarouselTrack(c);if(!track)return;
			var count=tcCarouselCount(track);
			idx=Math.max(0,Math.min(idx,count-1));
			track.scrollTo({left:idx*track.offsetWidth,behavior:'smooth'});
			}
			document.addEventListener('click',function(e){
			var next=e.target.closest('[data-tc-carousel-next]');
			var prev=e.target.closest('[data-tc-carousel-prev]');
			var dot=e.target.closest('[data-tc-carousel-dot]');
			if(!next&&!prev&&!dot)return;
			var c=e.target.closest('[data-tc-carousel]');if(!c)return;
			var track=tcCarouselTrack(c);if(!track)return;
			var idx=tcCarouselIndex(track);
			if(next)tcCarouselGo(c,idx+1);
			if(prev)tcCarouselGo(c,idx-1);
			if(dot)tcCarouselGo(c,parseInt(dot.getAttribute('data-tc-carousel-dot'),10));
			});
			document.addEventListener('keydown',function(e){
			var c=e.target.closest('[data-tc-carousel]');if(!c)return;
			if(e.key!=='ArrowLeft'&&e.key!=='ArrowRight'&&e.key!=='Home'&&e.key!=='End')return;
			e.preventDefault();
			var track=tcCarouselTrack(c);if(!track)return;
			var idx=tcCarouselIndex(track);
			var count=tcCarouselCount(track);
			var rtl=(c.closest('[dir]')||document.documentElement).getAttribute('dir')==='rtl';
			if(e.key==='Home')tcCarouselGo(c,0);
			else if(e.key==='End')tcCarouselGo(c,count-1);
			else if((!rtl&&e.key==='ArrowRight')||(rtl&&e.key==='ArrowLeft'))tcCarouselGo(c,idx+1);
			else tcCarouselGo(c,idx-1);
			});
			function tcBindTracks(){
			document.querySelectorAll('[data-tc-carousel-track]').forEach(function(track){
			if(track.dataset.tcBound)return;track.dataset.tcBound='1';
			if('onscrollend' in window){
			track.addEventListener('scrollend',function(){tcUpdateDots(track);});
			}else{
			var t;track.addEventListener('scroll',function(){clearTimeout(t);t=setTimeout(function(){tcUpdateDots(track);},150);});
			}
			});
			}
			tcBindTracks();
			document.body.addEventListener('htmx:afterSettle',tcBindTracks);
			})();`
