(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[931],{5561:function(e,t,s){Promise.resolve().then(s.bind(s,8434)),Promise.resolve().then(s.bind(s,2010))},8434:function(e,t,s){"use strict";s.r(t),s.d(t,{default:function(){return u}});var n=s(7437),r=s(2265),o=s(3859),a=s(3167),i=s(1367);function l(){for(var e=arguments.length,t=Array(e),s=0;s<e;s++)t[s]=arguments[s];return(0,i.m6)((0,a.W)(t))}let c=e=>{let{email:t,subject:s="",body:r="",className:o,children:a}=e,i=s||r?"?":"";s&&(i+="subject=".concat(encodeURIComponent(s))),r&&(i+="".concat(s?"&":"","body=").concat(encodeURIComponent(r)));let c=l("no-underline",o||"");return(0,n.jsx)("a",{className:c,href:"mailto:".concat(t).concat(i),children:a})};var u=()=>{let e=(0,r.useRef)(null),t=(0,r.useRef)(null),s=(0,r.useRef)(null),a={idle:"idling",success:"success",failure:"failure"},[i,u]=(0,r.useState)({state:a.idle,msg:""});(0,r.useEffect)(()=>{i&&(s.current.disabled=!0)},[i]);let d=async()=>{var s,n;let r={email:null===(s=e.current)||void 0===s?void 0:s.value,csrf_token:null===(n=t.current)||void 0===n?void 0:n.value},o=window.location.origin;try{let e=await fetch(o,{method:"POST",credentials:"same-origin",headers:{"Content-Type":"application/json"},redirect:"follow",body:JSON.stringify(r)}),t=await e.json();console.log(t,"=======");let s=a.failure;t.success&&(s=a.success),u({state:s,msg:t.msg})}catch(e){u({state:status,msg:"Ah man, The Server is down."})}};return(0,n.jsxs)("section",{className:"w-full min-h-96 mx-8 bg-slate-950 flex flex-col lg:flex-row gap-4 px-4 pt-32 pb-16 justify-evenly",children:[(0,n.jsxs)("header",{className:"",children:[(0,n.jsx)("h2",{className:"text-7xl font-secondary capitalize text-[#eb4fa0]",children:"Let's talk"}),(0,n.jsxs)("div",{className:"font-primary text-white/[0.5] overflow-hidden mt-4",children:["We would love to work with you, help you building",(0,n.jsx)("br",{}),"your awesome product."]}),(0,n.jsxs)("p",{className:"text-purple-600 mt-8",children:[(0,n.jsx)(o.Z,{size:18,className:"text-purple-600 inline-block"}),(0,n.jsx)(c,{className:"px-2 text-purple-700",email:"amitava.dev@proton.me",subject:"Discussing Project Opportunities",body:"Hi!, we would like to discuss an opportunity with you.",children:"amitava.dev@proton.me"})]})]}),(0,n.jsxs)("div",{className:"w-1/2 flex flex-col",children:[(0,n.jsxs)("form",{className:"w-full",onSubmit:d,method:"POST",action:"/tbc/emails/leads#contact-us",id:"contact-us",children:[(0,n.jsx)("input",{type:"hidden",ref:t,name:"csrf_token",value:"{{.csrf_token}}"}),(0,n.jsxs)("div",{className:"mb-6",children:[(0,n.jsx)("label",{className:"block mb-2 text-sm font-medium text-gray-900 dark:text-white",children:"Email address"}),(0,n.jsx)("input",{name:"email",ref:e,type:"email",id:"email",className:"bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500",placeholder:"john.doe@company.com",required:!0})]}),(0,n.jsx)("button",{ref:s,type:"submit",className:"text-white bg-purple-700 hover:bg-purple-800 focus:ring-4 focus:outline-none focus:ring-blue-700 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center",children:"Submit"})]}),(0,n.jsx)("span",{className:l("inline-block max-w-[1/2] mt-4 text-white font-primary",i.state==a.idle?"hidden":i.state==a.success?"text-green-700":"text-red-600"),children:i.msg})]})]})}},2010:function(e,t,s){"use strict";s.r(t),s.d(t,{default:function(){return b}});var n=s(7437),r=s(2265),o=s(9631),a=s(8648),i=s(2311),l=s(4620),c=s(5174),u=s(9088);let d=e=>{var t,s;e=e>1?1:e<0?0:e;let n=null===(t="132b8d".match(/.{1,2}/g))||void 0===t?void 0:t.map(t=>parseInt(t,16)*(1-e)),r=null===(s="fd12d1".match(/.{1,2}/g))||void 0===s?void 0:s.map(t=>parseInt(t,16)*e);if(!n||!r)return"#{LEFT_COLOR}";let o=[0,1,2].map(e=>Math.min(Math.round(n[e]+r[e]),255)).reduce((e,t)=>(e<<8)+t,0).toString(16).padStart(6,"0");return"#".concat(o)},m=e=>d((e+15)/30),h=(e,t)=>Math.random()*(t-e)+e,p=Array.from({length:2500},(e,t)=>t+1).map(e=>{let t=h(7.5,15),s=Math.random()*Math.PI*2,n=Math.cos(s)*t,r=h(-2,2),o=m(n);return{idx:e,position:[n,Math.sin(s)*t,r],color:o}}),f=Array.from({length:2500/3},(e,t)=>t+1).map(e=>{let t=h(3.75,30),s=Math.random()*Math.PI*2,n=Math.cos(s)*t,r=h(-20,20),o=m(n);return{idx:e,position:[n,Math.sin(s)*t,r],color:o}}),x=e=>{let t=(0,r.useRef)(null);return(0,a.C)(s=>{var n;let{clock:r}=s;(null===(n=t.current)||void 0===n?void 0:n.rotation)&&e.inView?(console.log("rendering rotation"),t.current.rotation.z=.01*r.getElapsedTime(),t.current.rotation.x=-.0006*r.getElapsedTime()):console.log("rotation stopped")}),(0,n.jsxs)("group",{ref:t,children:[p.map(e=>(0,n.jsx)(g,{position:e.position,color:e.color},e.idx)),f.map(e=>(0,n.jsx)(g,{position:e.position,color:e.color},e.idx))]})},g=e=>{let{position:t,color:s}=e;return(0,n.jsx)(i.aL,{position:t,args:[.03,10,10],children:(0,n.jsx)("meshStandardMaterial",{emissive:s,emissiveIntensity:.5,roughness:.9,color:s})})};var b=()=>{let e=(0,r.useRef)(null),[t,s]=(0,r.useState)(!1),a=e=>()=>{s(e.current.getBoundingClientRect().y<0)};return(0,r.useEffect)(()=>{if(!e)return;let t=a(e);return window.addEventListener("scroll",t),()=>{window.removeEventListener("scroll",t)}}),(0,n.jsxs)("div",{className:"relative w-screen h-screen max-h-screen",children:[(0,n.jsxs)(o.Xz,{camera:{position:[-10,10.5,-10]},className:"bg-slate-950",children:[(0,n.jsx)("directionalLight",{}),(0,n.jsx)("pointLight",{position:[-30,0,-30],power:10}),(0,n.jsx)(x,{inView:!t}),(0,n.jsxs)(l.x,{children:[(0,n.jsx)(c.y,{focusDistance:0,focalLength:.05,bokehScale:1,height:480}),(0,n.jsx)(u.d,{luminanceThreshold:0,luminanceSmoothing:.9,height:300,opacity:3})]})]}),(0,n.jsx)("span",{ref:e,className:"troll-spanner"})]})}}},function(e){e.O(0,[689,723,796,971,69,744],function(){return e(e.s=5561)}),_N_E=e.O()}]);