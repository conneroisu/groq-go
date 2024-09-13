library IEEE;

use IEEE.std_logic_1164.all;

entity TPU_MV_Element is

  port(
    iCLK : in  std_logic;
    iX   : in  integer;
    iW   : in  integer;
    iLdW : in  integer;
    iY   : in  integer;
    oY   : out integer;
    oX   : out integer
    );

end TPU_MV_Element;

architecture structure of TPU_MV_Element is

  component Adder
    port(
         iCLK : in  std_logic;
         iA   : in  integer;
         iB   : in  integer;
         oC   : out integer
       );
  end component;
  
  component Multiplier
    port( 
         iCLK : in  std_logic;
         iA   : in  integer;
         iB   : in  integer;
         oP   : out integer
       );
  end component;
  
  component Reg
    port(
         iCLK : in  std_logic;
         iD   : in  integer;
         oQ   : out integer
       );
  end component;
  
  component RegLd
    port(
         iCLK : in  std_logic;
         iD   : in  integer;
         iLd  : in  integer;
         oQ   : out integer
       );
  end component;

  signal s_W   : integer;
  signal s_X1  : integer;
  signal s_Y1  : integer;
  signal s_WxX : integer;               -- Signal to carry stored W*X

begin
  -- Level 0: Conditionally load new W
  g_Weight : RegLd
    port map(
             iCLK => iCLK,
             iD   => iW,
             iLd  => iLdW,
             oQ   => s_W
           );
  -- Level 1: Delay X and Y, calculate W*X
  g_Delay1 : Reg
    port map(
             iCLK => iCLK,
             iD   => iX,
             oQ   => s_X1
           );
  g_Delay2 : Reg
    port map(
             iCLK => iCLK,
             iD   => iY,
             oQ   => s_Y1
           );
  g_Mult1 : Multiplier
    port map(
             iCLK => iCLK,
             iA   => iX,
             iB   => s_W,
             oP   => s_WxX
           );
  -- Level 2: Delay X, calculate Y += W*X
  g_Delay3 : Reg
    port map(
             iCLK => iCLK,
             iD   => s_X1,
             oQ   => oX
           );
  g_Add1 : Adder
    port map(
             iCLK => iCLK,
             iA   => s_WxX,
             iB   => s_Y1,
             oC   => oY
           );
    
end structure;
