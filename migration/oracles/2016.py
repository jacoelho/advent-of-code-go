"""Independent 2016 calculations; pass a checkout of the pinned Clojure source."""
import sys
from pathlib import Path
import re,json,hashlib,collections,itertools,functools,math,time
ROOT=Path(sys.argv[1])/'resources'/'2016';OUT={}
INLINE={5:'cxdnnyjw',13:'1362',14:'ihaygndm',16:'11011110011011101',17:'mmsxrhfx',19:'3004953'}
SOURCE={1:(242,150),2:('74921','A6B35'),3:(1050,1921),4:(137896,501),5:('f77a0e6e','999828ec'),6:('qqqluigu','lsoypmia'),7:(115,231),8:(119,'ZFHFSFOGPO'),9:(150914,11052855125),10:(116,23903),11:(33,57),12:(318077,9227731),13:(82,138),14:(15035,19968),15:(376777,3903937),16:('00000100100001100','00011010100010010'),17:('RLDUDRDDRR',590),18:(2013,20006289),19:(1815603,1410630),20:(23923783,125),21:('gbhafcde','bcfaegdh'),22:(1043,185),23:(10880,479007440),24:(430,700),25:(198,None)}
def raw(d):return INLINE[d] if d in INLINE else (ROOT/f'day{d:02}.txt').read_text()
def save(d,a,b=None,provenance='independent Python calculation',**extra):
 p=ROOT/f'day{d:02}.txt';OUT[str(d)]={'part1':str(a),'part2':str(b) if b is not None else None,'source_expected':[str(x) if x is not None else None for x in SOURCE[d]],'source_type':'comment only' if d in (11,23,25) else 'source test assertions','input_path':str(p) if d not in INLINE else None,'inline_input':INLINE.get(d),'input_sha256':hashlib.sha256(raw(d).encode()).hexdigest(),'provenance':provenance,**extra};print(d,a,b,file=sys.stderr,flush=True)
# Fast independent oracles; heavy hashing and VMs are separate serial jobs.
p=0j;direction=1j;seen={p};first=None
for turn,n in re.findall(r'([LR])(\d+)',raw(1)):
 direction*=1j if turn=='L' else -1j
 for _ in range(int(n)):
  p+=direction
  if p in seen and first is None:first=int(abs(p.real)+abs(p.imag))
  seen.add(p)
save(1,int(abs(p.real)+abs(p.imag)),first)
def keypad(rows):
 keys={(x,y):c for y,row in enumerate(rows) for x,c in enumerate(row) if c!=' '};p=next(p for p,c in keys.items() if c=='5');code=''
 for line in raw(2).splitlines():
  for c in line:
   dx,dy={'U':(0,-1),'D':(0,1),'L':(-1,0),'R':(1,0)}[c];q=p[0]+dx,p[1]+dy
   if q in keys:p=q
  code+=keys[p]
 return code
save(2,keypad(['123','456','789']),keypad(['  1  ',' 234 ','56789',' ABC ','  D  ']))
tri=[list(map(int,l.split())) for l in raw(3).splitlines()];valid=lambda a:sum(a)>2*max(a);columns=[v for i in range(0,len(tri),3) for v in zip(*tri[i:i+3])];save(3,sum(map(valid,tri)),sum(map(valid,columns)))
sector=0;target=None
for l in raw(4).splitlines():
 name,num,check=re.fullmatch(r'([a-z-]+)-(\d+)\[([a-z]+)\]',l).groups();freq=collections.Counter(name.replace('-',''));real=''.join(sorted(freq,key=lambda c:(-freq[c],c))[:5])==check
 if real:
  sector+=int(num);decoded=''.join(' ' if c=='-' else chr((ord(c)-97+int(num))%26+97) for c in name)
  if decoded=='northpole object storage':target=int(num)
save(4,sector,target)
lines=raw(6).splitlines();cols=[collections.Counter(c) for c in zip(*lines)];save(6,''.join(min(c,key=lambda x:(-c[x],x)) for c in cols),''.join(min(c,key=lambda x:(c[x],x)) for c in cols))
def abba(s):return any(a==d and b==c and a!=b for a,b,c,d in zip(s,s[1:],s[2:],s[3:]))
a=b=0
for l in raw(7).splitlines():
 parts=re.split(r'[\[\]]',l);outer=parts[::2];inner=parts[1::2];a+=any(map(abba,outer)) and not any(map(abba,inner));patterns={y+x+y for s in outer for x,y,z in zip(s,s[1:],s[2:]) if x==z and x!=y};b+=any(p in s for p in patterns for s in inner)
save(7,a,b)
grid=[[False]*50 for _ in range(6)]
for line in raw(8).splitlines():
 vals=list(map(int,re.findall(r'\d+',line)));a,b=vals
 if line.startswith('rect'):
  for y in range(b):
   for x in range(a):grid[y][x]=True
 elif 'row' in line:
  row=grid[a][:]
  for x,c in enumerate(row):grid[a][(x+b)%50]=c
 else:
  col=[row[a] for row in grid]
  for y,c in enumerate(col):grid[(y+b)%6][a]=c
rows=[''.join('#' if c else ' ' for c in row) for row in grid];save(8,sum(map(sum,grid)),'\n'.join(rows),rendered_rows=rows,human_reading='ZFHFSFOGPO');print('\n'.join(row.replace(' ','.') for row in rows),file=sys.stderr,flush=True)
s=raw(9).strip()
def decompressed(s,recurse):
 i=total=0
 while i<len(s):
  if s[i]!='(':total+=1;i+=1;continue
  close=s.index(')',i);size,times=map(int,s[i+1:close].split('x'));chunk=s[close+1:close+1+size];assert len(chunk)==size;total+=times*(decompressed(chunk,True) if recurse else size);i=close+1+size
 return total
save(9,decompressed(s,False),decompressed(s,True))
bots=collections.defaultdict(list);rules={};outputs=collections.defaultdict(list)
for l in raw(10).splitlines():
 nums=list(map(int,re.findall(r'\d+',l)))
 if l.startswith('value'):v,bot=nums;bots[bot].append(v)
 else:
  bot,lo,hi=nums;kinds=re.findall(r'to (bot|output)',l);rules[bot]=[(kinds[0],lo),(kinds[1],hi)]
comparison=None
while True:
 ready=[bot for bot,chips in bots.items() if len(chips)==2]
 if not ready:break
 bot=ready[0];lo,hi=sorted(bots.pop(bot))
 if (lo,hi)==(17,61):comparison=bot
 for v,(kind,idx) in zip((lo,hi),rules[bot]):(bots if kind=='bot' else outputs)[idx].append(v)
save(10,comparison,math.prod(outputs[i][0] for i in (0,1,2)))
# BFS over canonical sorted generator/chip floor pairs, exact legal cargo moves.
def elevator(pairs):
 start=(0,tuple(sorted(pairs)));queue=collections.deque([(start,0)]);seen={start}
 def safe(ps):return all(g==c or all(other_g!=c for other_g,_ in ps) for g,c in ps)
 while queue:
  (floor,pairs),steps=queue.popleft()
  if all(p==(3,3) for p in pairs):return steps,len(seen)
  cargo=[(i,t) for i,p in enumerate(pairs) for t,f in enumerate(p) if f==floor]
  for dest in (floor+1,floor-1):
   if not 0<=dest<4:continue
   if dest<floor and not any(f<floor for p in pairs for f in p):continue
   for count in (1,2):
    for moved in itertools.combinations(cargo,count):
     ps=[list(p) for p in pairs]
     for i,t in moved:ps[i][t]=dest
     if not safe(ps):continue
     nxt=(dest,tuple(sorted(map(tuple,ps))))
     if nxt not in seen:seen.add(nxt);queue.append((nxt,steps+1))
 raise ValueError('no solution')
pairs=[(0,0)]+[(1,2)]*4;a,sa=elevator(pairs);b,sb=elevator(pairs+[(0,0)]*2);save(11,a,b,states_seen=[sa,sb])
# Explicit maze BFS; do not stop until both distance and50-step set known.
fav=int(raw(13));q=collections.deque([((1,1),0)]);seen={(1,1)};dist=None;reachable=0
while q:
 (x,y),d=q.popleft()
 if d<=50:reachable+=1
 if (x,y)==(31,39):dist=d
 if dist is not None and d>=max(50,dist):break
 for xx,yy in ((x+1,y),(x-1,y),(x,y+1),(x,y-1)):
  if xx<0 or yy<0 or (xx,yy) in seen:continue
  if (xx*xx+3*xx+2*xx*yy+yy+yy*yy+fav).bit_count()%2:continue
  seen.add((xx,yy));q.append(((xx,yy),d+1))
save(13,dist,reachable)
# CRT merge for disc release time.
discs=[tuple(map(int,re.fullmatch(r'Disc #(\d+) has (\d+) positions; at time=0, it is at position (\d+)\.',l).groups())) for l in raw(15).splitlines()]
def release(ds):
 t=0;mod=1
 for number,n,start in ds:
  g=math.gcd(mod,n);need=-number-start-t;assert need%g==0;t+=mod*((need//g*pow(mod//g,-1,n//g))%(n//g));mod=math.lcm(mod,n);t%=mod
 return t
save(15,release(discs),release(discs+[(len(discs)+1,11,0)]))
# Dragon expansion in bytes, block parity gives final checksum without full reductions.
seed=raw(16).encode();a=None
for size in (272,35651584):
 data=bytearray(seed)
 while len(data)<size:data+=b'0'+data[::-1].translate(bytes.maketrans(b'01',b'10'))
 del data[size:];block=size&-size;checksum=''.join('1' if data[i:i+block].count(49)%2==0 else '0' for i in range(0,size,block))
 if size==272:a=checksum
 else:save(16,a,checksum)
del data
# Path-dependent door search; vault goal states are terminal.
salt=raw(17).encode();queue=collections.deque([(0,0,'')]);short=None;longest=0
while queue:
 x,y,path=queue.popleft()
 if (x,y)==(3,3):
  if short is None:short=path
  longest=max(longest,len(path));continue
 h=hashlib.md5(salt+path.encode()).hexdigest()
 for ch,dx,dy,door in zip('UDLR',[0,0,-1,1],[-1,1,0,0],h[:4]):
  xx,yy=x+dx,y+dy
  if door in 'bcdef' and 0<=xx<4 and 0<=yy<4:queue.append((xx,yy,path+ch))
save(17,short,longest)
s=raw(18).strip();bits=sum((c=='^')<<i for i,c in enumerate(s));mask=(1<<len(s))-1;safe=0;a=None
for row in range(1,400001):
 safe+=len(s)-bits.bit_count();bits=((bits<<1)^(bits>>1))&mask
 if row==40:a=safe
save(18,a,safe)
n=int(raw(19));power2=1<<(n.bit_length()-1);power3=1
while power3*3<=n:power3*=3
save(19,2*(n-power2)+1,n if n==power3 else n-power3 if n<=2*power3 else 2*n-3*power3)
# Inclusive uint32 complement with gaps before, between and after intervals.
ranges=sorted(tuple(map(int,l.split('-'))) for l in raw(20).splitlines());current=0;first=None;count=0
for lo,hi in ranges:
 if lo>current:
  if first is None:first=current
  count+=lo-current
 current=max(current,hi+1)
if current<2**32:
 if first is None:first=current
 count+=2**32-current
save(20,first,count)
# Password transform plus exhaustive inverse over8! permutations (not inverse code).
ops=raw(21).splitlines()
def scramble(s):
 s=list(s)
 for l in ops:
  p=l.split()
  if p[0]=='swap':
   a,b=(int(p[2]),int(p[5])) if p[1]=='position' else (s.index(p[2]),s.index(p[5]));s[a],s[b]=s[b],s[a]
  elif p[0]=='rotate':
   n=(s.index(p[-1])+1+(s.index(p[-1])>=4)) if p[1]=='based' else int(p[2])*(1 if p[1]=='right' else -1);n%=len(s);s=s[-n:]+s[:-n] if n else s
  elif p[0]=='reverse':a,b=int(p[2]),int(p[4]);s[a:b+1]=reversed(s[a:b+1])
  else:a,b=int(p[2]),int(p[5]);s.insert(b,s.pop(a))
 return ''.join(s)
a=scramble('abcdefgh');matches=[''.join(p) for p in itertools.permutations('abcdefgh') if scramble(p)=='fbgdceah'];assert len(matches)==1;save(21,a,matches[0])
# Storage proof and exact state compression under verified capacity inequalities.
nodes={}
for l in raw(22).splitlines():
 m=re.search(r'node-x(\d+)-y(\d+)\s+(\d+)T\s+(\d+)T\s+(\d+)T',l)
 if m:x,y,size,used,avail=map(int,m.groups());nodes[x,y]=(size,used,avail)
a=sum(used>0 and p!=q and used<=avail for p,(_,used,_) in nodes.items() for q,(_,_,avail) in nodes.items());empty=next(p for p,v in nodes.items() if v[1]==0);movable={p for p,v in nodes.items() if v[1]<=nodes[empty][0]};walls=set(nodes)-movable;mincap=min(nodes[p][0] for p in movable);maxcap=max(nodes[p][0] for p in movable);payloads=[nodes[p][1] for p in movable if nodes[p][1]>0];minpayload=min(payloads);maxpayload=max(payloads);wallminused=min(nodes[p][1] for p in walls);wallmaxavail=max(nodes[p][2] for p in walls)
assert maxpayload<=mincap and 2*minpayload>maxcap and wallminused>maxcap and wallmaxavail<minpayload
# Wall-to-wall transfers also impossible.
assert all(nodes[p][1]>nodes[q][2] for p in walls for q in walls if p!=q)
goal=(max(x for x,y in nodes if y==0),0);q=collections.deque([((empty,goal),0)]);seen={(empty,goal)}
while q:
 (emp,g),d=q.popleft()
 if g==(0,0):b=d;break
 x,y=emp
 for p in ((x+1,y),(x-1,y),(x,y+1),(x,y-1)):
  if p not in movable:continue
  nxt=(p,emp if p==g else g)
  if nxt not in seen:seen.add(nxt);q.append((nxt,d+1))
save(22,a,b,compression_proof=dict(mincap=mincap,maxcap=maxcap,minpayload=minpayload,maxpayload=maxpayload,wallminused=wallminused,wallmaxavail=wallmaxavail),states_seen=len(seen))
# Maze all-pairs BFS and exhaustive point-order tour.
grid=raw(24).splitlines();points={c:(x,y) for y,row in enumerate(grid) for x,c in enumerate(row) if c.isdigit()};dist={}
for label,start in points.items():
 q=collections.deque([(start,0)]);seen={start}
 while q:
  (x,y),d=q.popleft()
  if grid[y][x].isdigit():dist[label,grid[y][x]]=d
  for xx,yy in ((x+1,y),(x-1,y),(x,y+1),(x,y-1)):
   if 0<=yy<len(grid) and 0<=xx<len(grid[yy]) and grid[yy][xx]!='#' and (xx,yy) not in seen:seen.add((xx,yy));q.append(((xx,yy),d+1))
costs=[]
for p in itertools.permutations(set(points)-{'0'}):
 path=('0',)+p;cost=sum(dist[a,b] for a,b in zip(path,path[1:]));costs.append((cost,cost+dist[path[-1],'0']))
save(24,min(a for a,b in costs),min(b for a,b in costs))
# Serial heavy-hash continuation, invoked explicitly rather than together with VM jobs.
def hash_password(salt,seconds=70):
 started=time.monotonic();key=salt.encode();plain='';position=[None]*8;last=None
 for i in itertools.count():
  if i%100000==0 and time.monotonic()-started>seconds:raise TimeoutError('password hash budget exceeded')
  h=hashlib.md5(key+str(i).encode()).digest()
  if h[0] or h[1] or h[2]>=16:continue
  if len(plain)<8:plain+=format(h[2],'x')
  p=h[2]
  if p<8 and position[p] is None:position[p]=format(h[3]>>4,'x')
  if len(plain)==8 and all(c is not None for c in position):return plain,''.join(position),i

def pad_keys(salt,stretch,count=64,seconds=70):
 started=time.monotonic();features=[];key=salt.encode();keys=[]
 def feature(i):
  s=key+str(i).encode()
  for _ in range(2017 if stretch else 1):s=hashlib.md5(s).hexdigest().encode()
  text=s.decode();m=re.search(r'(.)\1\1',text);return m[1] if m else None,{text[j] for j in range(28) if len(set(text[j:j+5]))==1}
 for i in itertools.count():
  if i%100==0 and time.monotonic()-started>seconds:raise TimeoutError('pad key hash budget exceeded')
  while len(features)<=i+1000:features.append(feature(len(features)))
  triple,_=features[i]
  if triple is not None and any(triple in quint for _,quint in features[i+1:i+1001]):keys.append(i)
  if len(keys)==count:return keys
# Independent guarded accelerated VM; validate acceleration against literal traces for manageable seeds.
def vm(program,initial,accelerate=True,clock=False,limit=20000000):
 code=[l.split() for l in program.splitlines()];r=dict.fromkeys('abcd',0);r.update(initial);ip=0;steps=0;expected=0;seen={};outputs=0;toggles=[]
 def val(s):return r[s] if s in r else int(s)
 while 0<=ip<len(code):
  steps+=1
  if steps>limit:raise TimeoutError('VM step budget exceeded')
  if clock:
   state=(ip,tuple(r[k] for k in 'abcd'),expected,tuple(map(tuple,code)))
   if state in seen:
    period=outputs-seen[state]
    return period>0,{'outputs':outputs,'output_period':period,'repeated_state':[ip,[r[k] for k in 'abcd'],expected],'steps':steps}
   seen[state]=outputs
  if accelerate and ip+5<len(code):
   b=code[ip:ip+6]
   if [x[0] for x in b]==['cpy','inc','dec','jnz','dec','jnz']:
    src,c=b[0][1:];a=b[1][1];d=b[4][1]
    if all(k in r for k in (a,c,d)) and len({a,c,d})==3 and src not in (a,c,d) and b[2]==['dec',c] and b[3]==['jnz',c,'-2'] and b[5]==['jnz',d,'-5'] and val(src)>0 and r[d]>0:
     r[a]+=val(src)*r[d];r[c]=r[d]=0;ip+=6;continue
  if accelerate and ip+2<len(code):
   b=code[ip:ip+3]
   if len(b[0])==len(b[1])==2 and b[0][0]=='inc' and b[1][0]=='dec':
    a=b[0][1];c=b[1][1]
    if a in r and c in r and a!=c and b[2]==['jnz',c,'-2'] and r[c]>0:r[a]+=r[c];r[c]=0;ip+=3;continue
  ins=code[ip];op=ins[0];args=ins[1:]
  if op=='cpy':
   if args[1] in r:r[args[1]]=val(args[0])
  elif op=='inc':
   if args[0] in r:r[args[0]]+=1
  elif op=='dec':
   if args[0] in r:r[args[0]]-=1
  elif op=='jnz':
   if val(args[0])!=0:ip+=val(args[1]);continue
  elif op=='tgl':
   target=ip+val(args[0])
   if 0<=target<len(code):
    before=code[target][0];code[target][0]=('dec' if before=='inc' else 'inc') if len(code[target])==2 else ('cpy' if before=='jnz' else 'jnz');toggles.append([ip,target,before,code[target][0]])
  elif op=='out':
   out=val(args[0])
   if clock and out!=expected:return False,{'outputs':outputs,'wrong_output':out,'expected':expected,'steps':steps}
   expected^=1;outputs+=1
  else:raise ValueError(ins)
  ip+=1
 return (False,{'outputs':outputs,'halted':True}) if clock else (r,toggles,steps)

# Hash jobs run serially; each function enforces its own elapsed-time limit.
a,b,last=hash_password(raw(5));save(5,a,b,last_hash_index=last)
normal=pad_keys(raw(14),False);stretched=pad_keys(raw(14),True)
save(14,normal[-1],stretched[-1],normal_key_indices=normal,stretched_key_indices=stretched)
a,_,_=vm(raw(12),{});b,_,_=vm(raw(12),{'c':1});save(12,a['a'],b['a'])
crosschecks={}
for seed in (6,7,8):
 literal,_,steps=vm(raw(23),{'a':seed},accelerate=False)
 accelerated,_,faststeps=vm(raw(23),{'a':seed})
 assert literal==accelerated
 crosschecks[str(seed)]={'a':literal['a'],'literal_steps':steps,'accelerated_steps':faststeps}
a,_,_=vm(raw(23),{'a':7});b,toggles,steps=vm(raw(23),{'a':12})
save(23,a['a'],b['a'],vm_baseline_crosschecks=crosschecks,toggles_seed12=toggles,vm_steps_seed12=steps)
rejections=collections.Counter()
for candidate in range(10000):
 valid,proof=vm(raw(25),{'a':candidate},clock=True)
 if valid:
  save(25,candidate,clock_proof=proof,rejected_smaller_candidates=candidate,rejection_summary=dict(rejections))
  break
 assert 'wrong_output' in proof
 rejections['wrong_output']+=1
else:raise RuntimeError('no clock proof within candidate budget')
OUT['8']['part2']=OUT['8']['human_reading']
for number,result in OUT.items():
 assert [result['part1'],result['part2']]==result['source_expected'],number
assert len(OUT)==25
print(json.dumps(OUT,indent=2))
