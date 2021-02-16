-- Card class

local physics = require 'physics'
physics.start()
-- physics.setGravity(0, 0.1)  -- 9.8
physics.setGravity(0, 0)  -- 9.8
-- print('physics.engineVersion', physics.engineVersion)

local Card = {
  -- prototype object
  pack = nil,
  suit = nil,
  ordinal = nil,
  prone = true,
  id = nil,
  color = nil,
  imgNum = nil,
  grp = nil,
    shadowRect = nil, -- drop shadow only shown when card is moving
    frontRect = nil, -- the front face of the card
    backRect = nil, -- the back face of the card
  owner = nil,  -- the pile this card belongs to, assigned by Pile:push()
}
Card.__index = Card

function Card.new(o)
  o = o or {}
  setmetatable(o, Card)

  o.id = o:encodeId()
  -- assert(string.len(o.id)==4)

  if o.suit == 'Heart' or o.suit == 'Diamond' then
    o.color = 'Red'
  else
    o.color = 'Black'
  end

  -- images in retro imagesheet are in Club, Diamond, Heart, Spade order
  if o.suit == 'Club' then
    o.imgNum = o.ordinal
  elseif o.suit == 'Diamond' then
    o.imgNum = 13 + o.ordinal
  elseif o.suit == 'Heart' then
    o.imgNum = 26 + o.ordinal
  elseif o.suit == 'Spade' then
    o.imgNum = 39 + o.ordinal
  else
    trace('Unknown suit in Card constructor', o.suit)
  end

  o:_build()

  return o
end

function Card:__tostring()
  return self.id
end

function Card:destroy()
  if self.grp then
    self.grp:removeEventListener('tap', self)
    self.grp:removeEventListener('touch', self)
    display.remove(self.grp)
    self.grp = nil
  end
end

function Card:encodeId()
  return string.format('%u%s%02u', self.pack, string.sub(self.suit, 1, 1), self.ordinal)
end

--[[
function Card.decodeId(id)  --- static function
  local suits = {H='Heart', C='Club', D='Diamond', S='Spade'}
  local obj = {
    pack = tonumber(string.sub(id, 1, 1)),
    suit = suits[string.sub(id, 2, 2)],
    ordinal = tonumber(string.sub(id, 3, 4)),
  }
  return obj
end
]]

function Card:_build()

  self:destroy()

  self.grp = display.newGroup()
  self.grp.x = _G.gm.stock.x  --display.contentCenterX
  self.grp.y = _G.gm.stock.y  --display.contentCenterY

  local SHADOW = _G.CARDWIDTH / 20
  self.shadowRect = display.newRoundedRect(self.grp, SHADOW, SHADOW, _G.CARDWIDTH, _G.CARDHEIGHT, _G.CARDRADIUS)
  self.shadowRect:setFillColor(unpack(_G.OPSOLE_COLORS.black))
  self.shadowRect.alpha = 0.33
  self.shadowRect.isVisible = false

  if _G.OPSOLE_SETTINGS.retroCards and _G.cardImageSheet then
    self.frontRect = display.newImage(self.grp, _G.cardImageSheet, self.imgNum)
    self.frontRect.width = _G.CARDWIDTH
    self.frontRect.height = _G.CARDHEIGHT
      -- kludge because retro .png does not have black border on cards
      local box = display.newRoundedRect(self.grp, 1, 1, _G.CARDWIDTH-2, _G.CARDHEIGHT-2, _G.CARDRADIUS-1)
      box:setFillColor(0,0)
      box:setStrokeColor(unpack(_G.OPSOLE_COLORS.black))
      box.strokeWidth = _G.STROKEWIDTH / 2
  else

    -- assert(self.color=='Red' or self.color=='Black')
    local fillColor = _G.OPSOLE_COLORS.black
    if self.color == 'Red' then
      fillColor = _G.OPSOLE_COLORS.red
    end

    self.frontRect = display.newRoundedRect(self.grp, 0, 0, _G.CARDWIDTH, _G.CARDHEIGHT, _G.CARDRADIUS)
    self.frontRect:setFillColor(unpack(_G.OPSOLE_COLORS.white))
    self.frontRect:setStrokeColor(unpack(_G.OPSOLE_COLORS.border))
    self.frontRect.strokeWidth = _G.STROKEWIDTH

    local cardFace = _G.gm.info.CardFace or 'Half'

    local faceValues = {'A','2','3','4','5','6','7','8','9','10','J','Q','K'}

    local ordx, ordy = -(_G.CARDWIDTH / 4), -(_G.CARDHEIGHT / 3)
    if self.ordinal == 10 then ordx = ordx + _G.CARDGAPX - 1 end
    self.ordText1 = display.newText(self.grp, faceValues[self.ordinal], ordx, ordy, _G.CARDFONT, _G.CARDWIDTH / 2)
    self.ordText1:setFillColor(unpack(fillColor))
    if cardFace == 'Full' then
      ordx, ordy = (_G.CARDWIDTH / 4), (_G.CARDHEIGHT / 3)
      if self.ordinal == 10 then ordx = ordx - _G.CARDGAPX - 1 end
      self.ordText2 = display.newText(self.grp, faceValues[self.ordinal], ordx, ordy, _G.CARDFONT, _G.CARDWIDTH / 2)
      self.ordText2:setFillColor(unpack(fillColor))
      self.ordText2:rotate(180)
    end

    local suitGlyphs = {Heart='♥', Club='♣', Spade='♠', Diamond='♦'}

    local suitx, suity = (_G.CARDWIDTH / 4), -(_G.CARDHEIGHT / 3)
    self.suitText1 = display.newText(self.grp, suitGlyphs[self.suit], suitx, suity, _G.CARDFONT, _G.CARDWIDTH / 3)
    self.suitText1:setFillColor(unpack(fillColor))
    if cardFace == 'Full' then
      suitx, suity = -(_G.CARDWIDTH / 4), (_G.CARDHEIGHT / 3)
      self.suitText2 = display.newText(self.grp, suitGlyphs[self.suit], suitx, suity, _G.CARDFONT, _G.CARDWIDTH / 3)
      self.suitText2:setFillColor(unpack(fillColor))
      self.suitText2:rotate(180)
    end

    if cardFace == 'Half' then
      if self.ordinal == 1 and self.suit == 'Spade' then
        -- original image is 280x180
        local width = _G.CARDWIDTH
        local height = width / (280/180)
        local raccoonImg = display.newImageRect(self.grp, _G.raccoon.filename, _G.raccoon.baseDir, width, height)
        raccoonImg.y = _G.CARDHEIGHT / 5
      elseif self.ordinal == 1 or self.ordinal > 10 then
        local bigsuit = display.newText(self.grp, suitGlyphs[self.suit], 0, _G.CARDHEIGHT / 5, _G.CARDFONT, _G.CARDWIDTH * 0.666)
        bigsuit:setFillColor(unpack(fillColor))
      end
    end

  end

  if _G.OPSOLE_SETTINGS.retroCards then
    self.backRect = display.newImage(self.grp, _G.cardbackImageSheet, _G.cardbackImageSheetIndex)
    self.backRect.width = _G.CARDWIDTH
    self.backRect.height = _G.CARDHEIGHT
  else
    self.backRect = display.newRoundedRect(self.grp, 0, 0, _G.CARDWIDTH, _G.CARDHEIGHT, _G.CARDRADIUS)
    self.backRect:setFillColor(unpack(_G.OPSOLE_CARDCOLORS[_G.OPSOLE_SETTINGS.cardColor]))
    self.backRect:setStrokeColor(unpack(_G.OPSOLE_COLORS.border))
    self.backRect.strokeWidth = _G.STROKEWIDTH
  end

  -- cards are make (in the Stock) face down
  self.prone = true
  self.backRect.isVisible = true

  -- assert(self.grp[1]==self.shadowRect)
  -- assert(self.grp[2]==self.frontRect)

  self.grp:addEventListener('tap', self)
  if _G.gm.info.Touch ~= 'Off' then
    self.grp:addEventListener('touch', self)
  end

  _G.OPSOLE_GROUPS.cards:insert(self.grp)

end

function Card:resetColor()
  if not _G.OPSOLE_SETTINGS.retroCards then
    self.backRect:setFillColor(unpack(_G.OPSOLE_CARDCOLORS[_G.OPSOLE_SETTINGS.cardColor]))
  end
end

function Card:getRect()
  -- currently, the img.x/y of the card is in it's center
  local halfX = _G.CARDWIDTH / 2
  local halfY = _G.CARDHEIGHT / 2
  return self.grp.x - halfX, self.grp.y - halfY, self.grp.x + halfX, self.grp.y + halfY
end

function Card:tap(event)
  -- dummy tap listener to pretend we handle this
  -- otherwise click on a stock card falls though to stock pile

  -- trace('Card:tap', self.id, event.name, event.numTaps, event.x, event.y)
  if _G.gm.info.Touch == 'Off' then
    -- let owner decide if this card should be toggled
    -- and what do do if it is or isn't
    self.owner:onClick(self)
  end

  return true
end

function Card:touch(event)

  -- event.target is self.grp, even if we clicked on titlebar.hamburger

  -- if event.target == _G.gm.titlebar.hamburger then
  --   trace('event.target is hamburger')
  -- elseif event.target == self then
  --   trace('event.target is self')
  -- elseif event.target == self.grp then
  --   trace('event.target is self.grp')
  -- end

  if event.phase == 'began' then

    if self.grabbedTail then trace('WARNING: touch began on a grabbed tail') end

    self.touchStartTime = event.time

    local src = self.owner
    if src:canGrab(self) then
      self.grabbedTail = src:getTail(self)
      for _,c in ipairs(self.grabbedTail) do
        c.originalX = c.grp.x
        c.originalY = c.grp.y
        c.offsetX = event.x - c.grp.x
        c.offsetY = event.y - c.grp.y

        c.shadowRect.isVisible = true

        c.grp:toFront()
      end
    else
      if not self.grp.getLinearVelocity then
        self:shake()
      end
    end

    -- "The ended phase is never triggered (and has never been) unless the touch ends on that object."
    display.getCurrentStage():setFocus(event.target)  -- stop orphaned cards

  elseif event.phase == 'moved' then

    if self.grabbedTail then
--[[
      do
        if not Util.inRect(event.x, event.y, self.grp.x - (_G.CARDWIDTH / 2), self.grp.y - (_G.CARDHEIGHT / 2), self.grp.x + (_G.CARDWIDTH / 2), self.grp.y + (_G.CARDHEIGHT / 2)) then
          trace('WARNING: touch no longer over card', tostring(self))
          -- for _,orphan in ipairs(self.grabbedTail) do
          --   orphan:transitionTo(orphan.originalX, orphan.originalY)
          -- end
        end
      end
]]
      local newx, newy = event.x - self.offsetX, event.y - self.offsetY
      if newx < 0 or newx > display.contentWidth or newy < 0 or newy > display.contentHeight then
        -- trace('WARNING: dragging a card off screen')
      else
        for _, c in ipairs(self.grabbedTail) do -- pairs() doesn't guarantee order (ipairs() does)
          c.grp:toFront()
          c.grp.x = event.x - c.offsetX
          c.grp.y = event.y - c.offsetY
        end
      end
    end

  elseif event.phase == 'ended' or event.phase == 'cancelled' then

    -- assert(event.xStart)
    -- assert(event.yStart)

    if self.grabbedTail then
      -- trace('touch ended or cancelled', event.xStart, event.yStart)
      -- https://www.quora.com/What-should-be-the-duration-of-a-simple-touch-tap-event-not-long-press-in-a-mobile-app
      if (event.time < self.touchStartTime + 125) or (event.x == event.xStart and event.y == event.yStart) then
        for _,c in ipairs(self.grabbedTail) do
          c.grp.x = c.originalX
          c.grp.y = c.originalY
        end
        -- We recommend you also include a target property in your event to the event so that your listener can know which object received the event.
        -- trace('faking a tap on', tostring(self), 'in pile', self.owner.info.Class)
        local tmp = self.grabbedTail  -- kludge of the year to detect if card is being clicked or dragged; see Waste:canAccept
        self.grabbedTail = nil
        self.owner:onClick(self)
        self.grabbedTail = tmp
      else
        -- trace('trying to get new owner')
        local dst = _G.gm:getNewOwner(self)
        if dst then
          _G.gm:moveCards(self.owner, dst, #self.grabbedTail)
        else
          -- trace('no owner for', tostring(self))
          for _,c in ipairs(self.grabbedTail) do
            c:transitionTo(c.originalX, c.originalY)
          end
        end
      end
      for _,c in ipairs(self.grabbedTail) do
        c.shadowRect.isVisible = false
      end
      self.grabbedTail = nil
    end

    if self.grp.getLinearVelocity then
      self:flip()
      self.grp:setLinearVelocity(math.random(-52,52), math.random(-52,52))
    end

    display.getCurrentStage():setFocus(nil)

  end
  return true
end

function Card:isProne()
  return self.prone
end

function Card:flip()
  local ms = _G.OPSOLE_SETTINGS.aniSpeed / 10
  self.prone = not self.prone -- immediately available, unlike self.backRect.isVisible
  transition.to(self.grp, {time=ms, xScale=0.1, onComplete=function()
    self.backRect.isVisible = not self.backRect.isVisible
    transition.to(self.grp, {time=ms, xScale=1})
  end})
end

function Card:transitionTo(x, y)
  transition.moveTo(self.grp, {x=x, y=y, time=_G.OPSOLE_SETTINGS.aniSpeed, transition=easing.outQuart})
  -- wobby move
  -- transition.moveTo(self.grp, {x=x, y=y, time=_G.OPSOLE_SETTINGS.aniSpeed, transition=easing.outQuart, onComplete=function() self:removePhysics() end})
  if y + (_G.CARDHEIGHT / 2) > display.contentHeight then
    trace('WARNING: card wandering off screen')
  end
  -- self.grp.x = x
  -- self.grp.y = y
end

function Card:shake()
  -- trace('shaking', tostring(self))
  transition.to(self.grp, {time=50, transition=easing.continuousLoop, x=self.grp.x + _G.CARDGAPX})
  transition.to(self.grp, {delay=50, time=50, transition=easing.continuousLoop, x=self.grp.x - _G.CARDGAPX})
end

function Card:markMovable(movable)
  -- self.grp is a group, the second object in it (self.grp[2]) is a card face rect
  self.frontRect:setFillColor(unpack(movable and _G.OPSOLE_COLORS.white or _G.OPSOLE_COLORS.offwhite))
end

-- function Card:setSelected(selected)
--   self.selected = selected
--   self.frontRect:setFillColor(unpack(selected and _G.OPSOLE_COLORS.selected or _G.OPSOLE_COLORS.offwhite))
-- end

-- function Card:toggleSelected()
--   self:setSelected(not self.selected)
-- end

function Card:getSaveableCard()
  return {id=self.id, prone=self:isProne()}
end

function Card:addPhysics()
  if self.grp.getLinearVelocity then
    return  -- already a physics object
  end
  -- the box is smaller than the card to allow some card overlap
  physics.addBody(self.grp, 'dynamic', { density=0.5, bounce=0.9, box={halfWidth=_G.CARDWIDTH/4, halfHeight=_G.CARDHEIGHT/4} } )
  self.grp:setLinearVelocity(math.random(-52,52), math.random(-52,52))
  -- self.grp.isFixedRotation = true
  self.shadowRect.isVisible = true
end

function Card:removePhysics()
  if not self.grp.getLinearVelocity then
    return  -- not a physics object
  end
  self.grp:rotate(-self.grp.rotation)
  physics.removeBody(self.grp)
  self.shadowRect.isVisible = false
end

function Card:constrain()
  if not self.grp.getLinearVelocity then
    return  -- not a physics object
  end
  local oldx, oldy = self.grp:getLinearVelocity()
  local newx, newy = oldx, oldy
  if self.grp.x > display.contentWidth then
    newx = math.random(-52, 0)
  elseif self.grp.x < 0 then
    newx = math.random(0, 52)
  end
  if self.grp.y > (display.contentHeight - _G.STATUSBARHEIGHT) then
    newy = math.random(-52, 0)
  elseif self.grp.y < _G.TITLEBARHEIGHT then
    newy = math.random(0, 52)
  end
  if newx ~= oldx or newy ~= oldy then
    self.grp:setLinearVelocity(newx, newy)
    -- self.grp.isFixedRotation = false
    if self:isProne() then self:flip() end
  end
end

return Card
