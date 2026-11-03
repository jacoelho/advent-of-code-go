"""Independent expectations; pass a checkout of the pinned Clojure repository."""
import sys
from pathlib import Path
import re,json,hashlib,itertools,functools,math,collections,heapq
ROOT=Path(sys.argv[1])/'resources'/'2015';OUT={}
SOURCE={1:(280,1797),2:(1588178,3783758),3:(2572,2631),4:(254575,1038736),5:(258,53),6:(377891,14110788),7:(3176,14710),8:(1333,2046),9:(207,804),10:(492982,6989950),11:('cqjxxyzz','cqkaabcc'),12:(191164,87842),13:(664,640),14:(2660,1256),15:(13882464,11171160),16:(103,405),17:(4372,4),18:(814,924),19:(576,207),20:(786240,831600),21:(111,188),22:(1269,1309),23:(255,334),24:(11266889531,77387711),25:(9132360,None)}
INLINE={4:'bgvyzdsv',10:'1321131112',11:'cqjxjnds',20:'34000000',22:'Hit Points: 58\nDamage: 9\n'}
def raw(d):return (ROOT/f'day{d:02}.txt').read_text() if d not in INLINE else INLINE[d]
def save(d,a,b=None,**kw):
 p=ROOT/f'day{d:02}.txt';OUT[str(d)]={'part1':str(a),'part2':str(b) if b is not None else None,'source_expected':[str(x) if x is not None else None for x in SOURCE[d]],'input_path':str(p) if d not in INLINE else None,'inline_input':INLINE.get(d),'input_sha256':hashlib.sha256(p.read_bytes()).hexdigest() if d not in INLINE else hashlib.sha256(INLINE[d].encode()).hexdigest(),'provenance':'independent Python oracle over pinned original fixture/inline seed',**kw};print(d,a,b,file=sys.stderr,flush=True)
# 1–3: direct character accounting; indexed boxes; set positions.
s=raw(1).strip();floor=0;first=None
for i,c in enumerate(s,1):
 floor+=1 if c=='(' else -1
 if floor<0 and first is None:first=i
save(1,floor,first)
a=b=0
for l in raw(2).splitlines():
 x,y,z=sorted(map(int,l.split('x')));a+=3*x*y+2*x*z+2*y*z;b+=2*(x+y)+x*y*z
save(2,a,b)
def visits(s):
 p=0j;seen={p}
 for c in s:p+={'^':1j,'v':-1j,'>':1,'<':-1}[c];seen.add(p)
 return seen
s=raw(3).strip();save(3,len(visits(s)),len(visits(s[::2])|visits(s[1::2])))
# MD5 minima: sequential enumeration proves no earlier qualifying nonce.
key=raw(4).encode();first=None
for n in itertools.count():
 h=hashlib.md5(key+str(n).encode()).digest()
 if h[0]==h[1]==0:
  if h[2]<16 and first is None:first=n
  if h[2]==0:save(4,first,n);break
# Nice strings use regex for independent predicates.
s=raw(5).splitlines();save(5,sum(sum(c in 'aeiou' for c in l)>=3 and re.search(r'(.)\1',l) is not None and re.search('ab|cd|pq|xy',l) is None for l in s),sum(re.search(r'(..).*\1',l) is not None and re.search(r'(.).\1',l) is not None for l in s))
# Light state: arbitrary-precision row bitsets for p1, separate scalar brightness.
from array import array
rows=[0]*1000;bright=array('I',[0])*1000000
for line in raw(6).splitlines():
 op,x,y,xx,yy=re.fullmatch(r'(turn on|turn off|toggle) (\d+),(\d+) through (\d+),(\d+)',line).groups();x,y,xx,yy=map(int,(x,y,xx,yy));mask=((1<<(xx-x+1))-1)<<x
 for j in range(y,yy+1):
  if op=='turn on':rows[j]|=mask
  elif op=='turn off':rows[j]&=~mask
  else:rows[j]^=mask
  for i in range(j*1000+x,j*1000+xx+1):bright[i]=max(0,bright[i]-1) if op=='turn off' else bright[i]+(1 if op=='turn on' else 2)
save(6,sum(row.bit_count() for row in rows),sum(bright));del bright
# Repeated topological circuit resolution, mask each result (no recursive Go algorithm).
rules={r:l.split() for l,r in (line.split(' -> ') for line in raw(7).splitlines())}
def circuit(override=None):
 vals={} if override is None else {'b':override}
 while 'a' not in vals:
  before=len(vals)
  for k,t in rules.items():
   if k in vals:continue
   args=[int(s) if s.isdigit() else vals.get(s) for s in t if s not in ('NOT','AND','OR','LSHIFT','RSHIFT')]
   if None in args:continue
   if len(t)==1:v=args[0]
   elif t[0]=='NOT':v=~args[0]
   elif t[1]=='AND':v=args[0]&args[1]
   elif t[1]=='OR':v=args[0]|args[1]
   elif t[1]=='LSHIFT':v=args[0]<<args[1]
   else:v=args[0]>>args[1]
   vals[k]=v&65535
  assert len(vals)>before
 return vals['a']
a=circuit();save(7,a,circuit(a))
# Python literal decoding independently interprets approved source string grammar.
import ast
lines=raw(8).splitlines();save(8,sum(len(s)-len(ast.literal_eval(s)) for s in lines),sum(2+s.count('"')+s.count('\\') for s in lines))
dist={};cities=set()
for l in raw(9).splitlines():
 a,b,d=re.fullmatch(r'(\w+) to (\w+) = (\d+)',l).groups();dist[a,b]=dist[b,a]=int(d);cities|={a,b}
cost=[sum(dist[a,b] for a,b in zip(p,p[1:])) for p in itertools.permutations(cities)];save(9,min(cost),max(cost))
s=raw(10);a=None
for phase in range(1,51):
 s=''.join(str(sum(1 for _ in g))+c for c,g in itertools.groupby(s))
 if phase==40:a=len(s)
save(10,a,len(s));del s
# Password search: distinct nonoverlapping pair letters.
def password(s):
 vals=list(map(ord,s))
 while True:
  i=7
  while i>=0 and vals[i]==122:vals[i]=97;i-=1
  assert i>=0;vals[i]+=1
  s=''.join(map(chr,vals))
  if any(c in s for c in 'iol'):continue
  if not any(vals[i]+2==vals[i+1]+1==vals[i+2] for i in range(6)):continue
  pairs=set(re.findall(r'(.)\1',s))
  if len(pairs)>=2:return s
a=password(raw(11));save(11,a,password(a))
data=json.loads(raw(12))
def total(v,drop):
 if isinstance(v,int):return v
 if isinstance(v,list):return sum(total(x,drop) for x in v)
 if isinstance(v,dict):return 0 if drop and 'red' in v.values() else sum(total(x,drop) for x in v.values())
 return 0
save(12,total(data,False),total(data,True))
# Circular happiness exhaustive arrangement of actual named participants.
happy={};people=set()
for l in raw(13).splitlines():
 a,sign,n,b=re.fullmatch(r'(\w+) would (gain|lose) (\d+) happiness units by sitting next to (\w+)\.',l).groups();happy[a,b]=int(n)*(1 if sign=='gain' else -1);people|={a,b}
def dinner(people):
 first=next(iter(people));return max(sum(happy.get((a,b),0)+happy.get((b,a),0) for a,b in zip((first,)+p,p+(first,))) for p in itertools.permutations(people-{first}))
save(13,dinner(people),dinner(people|{'me'}))
# Award every leader a point at every second, including final second.
deer=[tuple(map(int,re.findall(r'\d+',l))) for l in raw(14).splitlines()];scores=[0]*len(deer)
for second in range(1,2504):
 distances=[v*(second//(fly+rest)*fly+min(second%(fly+rest),fly)) for v,fly,rest in deer];best=max(distances)
 for i,d in enumerate(distances):scores[i]+=d==best
save(14,max(distances),max(scores))
# All nonnegative integer compositions of100, score after clamping attributes.
ing=[tuple(map(int,re.findall(r'-?\d+',l))) for l in raw(15).splitlines()]
def compositions(n,k):
 if k==1:yield(n,);return
 for i in range(n+1):
  for rest in compositions(n-i,k-1):yield(i,)+rest
best=best500=0
for portions in compositions(100,len(ing)):
 attrs=[sum(q*row[c] for q,row in zip(portions,ing)) for c in range(5)];score=math.prod(max(0,n) for n in attrs[:4]);best=max(best,score)
 if attrs[4]==500:best500=max(best500,score)
save(15,best,best500)
ref=dict(children=3,cats=7,samoyeds=2,pomeranians=3,akitas=0,vizslas=0,goldfish=5,trees=3,cars=2,perfumes=1);sues={}
for l in raw(16).splitlines():
 label,props=l.split(': ',1);sues[int(label.split()[1])]={k:int(v) for k,v in re.findall(r'(\w+): (\d+)',props)}
def aunt(part2):
 matches=[]
 for n,props in sues.items():
  if all((v>ref[k] if k in ('cats','trees') else v<ref[k] if k in ('pomeranians','goldfish') else v==ref[k]) if part2 else v==ref[k] for k,v in props.items()):matches.append(n)
 assert len(matches)==1;return matches[0]
save(16,aunt(False),aunt(True))
weights=list(map(int,raw(17).splitlines()));counts=collections.Counter({(0,0):1})
for w in weights:
 for (v,k),n in list(counts.items()):
  if v+w<=150:counts[v+w,k+1]+=n
matches={k:n for (v,k),n in counts.items() if v==150};save(17,sum(matches.values()),matches[min(matches)],minimum_container_count=min(matches),ways_by_count=matches)
# Bounded cellular automaton with forced corner seeds and every-step persistence.
grid=raw(18).splitlines();height=len(grid);width=len(grid[0]);corners={(0,0),(width-1,0),(0,height-1),(width-1,height-1)}
def lights(stuck):
 active={(x,y) for y,row in enumerate(grid) for x,c in enumerate(row) if c=='#'}
 if stuck:active|=corners
 for _ in range(100):
  counts=collections.Counter((x+dx,y+dy) for x,y in active for dx,dy in itertools.product((-1,0,1),repeat=2) if (dx or dy) and 0<=x+dx<width and 0<=y+dy<height);active={p for p,n in counts.items() if n==3 or n==2 and p in active}
  if stuck:active|=corners
 return len(active)
save(18,lights(False),lights(True))
# Molecule: enumerate single replacements on element tokens; prove lower bound via rank.
rules_s,molecule=raw(19).strip().split('\n\n');token=lambda s:tuple(re.findall(r'[A-Z][a-z]?|e',s));rules=[(token(a),token(b)) for a,b in (l.split(' => ') for l in rules_s.splitlines())];m=token(molecule);products=set()
for a,b in rules:
 for i in range(len(m)-len(a)+1):
  if m[i:i+len(a)]==a:products.add(m[:i]+b+m[i+len(a):])
def rank(t):return len(t)-t.count('Rn')-t.count('Ar')-2*t.count('Y')
assert all(rank(b)-rank(a)==1 for a,b in rules)
minimum=rank(m)-rank(('e',))
# Replay the independently discovered derivation; the rank bound proves minimality.
record=json.loads((Path(__file__).parents[1]/'2015-oracles.json').read_text())['days']['19']
witness=record['reverse_witness'];current=m
alphabet=set(m)|{t for a,b in rules for t in a+b}
assert all(not (lhs[0] in atom and lhs[0]!=atom) for lhs,_ in rules for atom in alphabet)
for index,before,after in witness:
 b,a=token(before),token(after)
 assert (a,b) in rules and current[index:index+len(b)]==b
 current=current[:index]+a+current[index+len(b):]
assert current==('e',) and len(witness)==minimum
save(19,len(products),minimum,rank_lower_bound=minimum,reverse_witness=witness)
# Sieve to a verified satisfying upper bound; every earlier house is fully tallied.
limit=1000000;target=34000000;gifts=array('I',[0])*(limit+1)
for elf in range(1,limit+1):
 for house in range(elf,limit+1,elf):gifts[house]+=10*elf
p1=next(h for h in range(1,limit+1) if gifts[h]>=target);gifts=array('I',[0])*(limit+1)
for elf in range(1,limit+1):
 for house in range(elf,min(limit,elf*50)+1,elf):gifts[house]+=11*elf
p2=next(h for h in range(1,limit+1) if gifts[h]>=target);save(20,p1,p2);del gifts
# Complete legal loadout enumeration and turn-count combat outcome.
weapons=[(8,4,0),(10,5,0),(25,6,0),(40,7,0),(74,8,0)];armors=[(0,0,0),(13,0,1),(31,0,2),(53,0,3),(75,0,4),(102,0,5)];rings=[(25,1,0),(50,2,0),(100,3,0),(20,0,1),(40,0,2),(80,0,3)];boss_hp,boss_damage,boss_armor=map(int,re.findall(r'\d+',raw(21)));wins=[];losses=[]
for w,a in itertools.product(weapons,armors):
 for count in range(3):
  for rs in itertools.combinations(rings,count):
   cost,damage,armor=[sum(item[i] for item in (w,a)+rs) for i in range(3)];player_turns=math.ceil(boss_hp/max(1,damage-boss_armor));boss_turns=math.ceil(100/max(1,boss_damage-armor));(wins if player_turns<=boss_turns else losses).append(cost)
save(21,min(wins),max(losses))
# Uniform-cost search at player-turn starts; expand an entire player/boss round.
def wizard(hard):
 start=(50,500,58,0,0,0);queue=[(0,start)];best={start:0};win=math.inf
 while queue:
  spent,state=heapq.heappop(queue)
  if spent>=win:break
  if best[state]!=spent:continue
  hp,mana,boss,shield,poison,recharge=state
  if hard:
   hp-=1
   if hp<=0:continue
  if poison:boss-=3
  if recharge:mana+=101
  shield=max(0,shield-1);poison=max(0,poison-1);recharge=max(0,recharge-1)
  if boss<=0:win=min(win,spent);continue
  for cost in (53,73,113,173,229):
   if cost>mana or cost==113 and shield or cost==173 and poison or cost==229 and recharge:continue
   hh,mm,bb,ss,pp,rr=hp,mana-cost,boss,shield,poison,recharge
   if cost==53:bb-=4
   if cost==73:bb-=2;hh+=2
   if cost==113:ss=6
   if cost==173:pp=6
   if cost==229:rr=5
   ns=spent+cost
   if bb<=0:win=min(win,ns);continue
   if pp:bb-=3
   if rr:mm+=101
   ss=max(0,ss-1);pp=max(0,pp-1);rr=max(0,rr-1)
   if bb<=0:win=min(win,ns);continue
   hh-=max(1,9-(7 if ss else 0))
   if hh<=0:continue
   nxt=(hh,mm,bb,ss,pp,rr)
   if ns<best.get(nxt,math.inf):best[nxt]=ns;heapq.heappush(queue,(ns,nxt))
 return win
save(22,wizard(False),wizard(True))
# Direct register-machine interpretation independent of Go dispatch.
program=raw(23).splitlines()
def vm(a):
 registers={'a':a,'b':0};ip=0;steps=0
 while 0<=ip<len(program):
  steps+=1;assert steps<1000000
  p=program[ip].replace(',','').split();op=p[0];r=p[1]
  if op=='hlf':registers[r]//=2
  if op=='tpl':registers[r]*=3
  if op=='inc':registers[r]+=1
  if op=='jmp':ip+=int(r);continue
  if op in ('jie','jio') and (registers[r]%2==0 if op=='jie' else registers[r]==1):ip+=int(p[2]);continue
  ip+=1
 return registers['b']
save(23,vm(0),vm(1))
# Validate remaining groups for every proposed smallest first group.
weights=tuple(sorted(map(int,raw(24).splitlines()),reverse=True))
def partition_qe(weights,groups):
 target=sum(weights)//groups;assert sum(weights)%groups==0;n=len(weights);full=(1<<n)-1
 @functools.cache
 def remainder(mask,count):
  if count==1:return True
  first=(mask&-mask).bit_length()-1
  def search(start,need,chosen):
   if need==0:return remainder(mask^chosen,count-1)
   for i in range(start,n):
    if mask>>i&1 and weights[i]<=need:
     if search(i+1,need-weights[i],chosen|1<<i):return True
   return False
  return search(first+1,target-weights[first],1<<first)
 for k in range(1,n+1):
  candidates=sorted((math.prod(weights[i] for i in ids),ids) for ids in itertools.combinations(range(n),k) if sum(weights[i] for i in ids)==target)
  for qe,ids in candidates:
   chosen=sum(1<<i for i in ids)
   if remainder(full^chosen,groups-1):return qe,k
 raise ValueError('no partition')
a,ka=partition_qe(weights,3);b,kb=partition_qe(weights,4);save(24,a,b,first_group_sizes=[ka,kb])
r,c=map(int,re.findall(r'\d+',raw(25)));diagonal=r+c-2;index=diagonal*(diagonal+1)//2+c;save(25,20151125*pow(252533,index-1,33554393)%33554393)
# Additional canonical checks requested by implementation worker.
def extra_canonical_checks():
 checks={};key=b'abcdef'
 for n in itertools.count():
  h=hashlib.md5(key+str(n).encode()).digest()
  if h[0]==h[1]==h[2]==0:checks['day04_abcdef_part2']=n;break
 s='1';forty=None
 for phase in range(1,51):
  s=''.join(str(sum(1 for _ in g))+c for c,g in itertools.groupby(s))
  if phase==40:forty=len(s)
 checks['day10_seed1_40_50']=[forty,len(s)]
 checks['day11_abcdefgh_part2']=password(password('abcdefgh'))
 print(json.dumps(checks),flush=True)
 return checks

print(json.dumps(OUT,indent=2))
